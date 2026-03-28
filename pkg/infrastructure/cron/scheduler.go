package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	cron_repository "github.com/make-bin/groundhog/pkg/domain/cron/repository"
	cron_service "github.com/make-bin/groundhog/pkg/domain/cron/service"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

const (
	defaultHeartbeatMs = 10000          // 10 seconds
	defaultCatchupMs   = 5 * 60 * 1000 // 5 minutes
)

// Scheduler is the background cron job scheduler.
// It periodically polls for due jobs and executes them asynchronously.
type Scheduler struct {
	cronRepo     cron_repository.CronRepository
	runLogRepo   cron_repository.CronRunLogRepository
	schedulerSvc cron_service.CronSchedulerService
	executor     *JobExecutor
	delivery     *DeliveryExecutor
	heartbeatMs  int64
	catchupMs    int64
	logger       logger.Logger

	mu      sync.Mutex
	running map[string]bool // jobID → currently running
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewScheduler creates a new Scheduler with the given dependencies.
func NewScheduler(
	cronRepo cron_repository.CronRepository,
	runLogRepo cron_repository.CronRunLogRepository,
	schedulerSvc cron_service.CronSchedulerService,
	executor *JobExecutor,
	delivery *DeliveryExecutor,
	heartbeatMs int64,
	log logger.Logger,
) *Scheduler {
	if heartbeatMs <= 0 {
		heartbeatMs = defaultHeartbeatMs
	}
	return &Scheduler{
		cronRepo:     cronRepo,
		runLogRepo:   runLogRepo,
		schedulerSvc: schedulerSvc,
		executor:     executor,
		delivery:     delivery,
		heartbeatMs:  heartbeatMs,
		catchupMs:    defaultCatchupMs,
		logger:       log,
		running:      make(map[string]bool),
	}
}

// Start begins the scheduler loop. It performs startup catchup logic and then
// enters the periodic tick loop.
func (s *Scheduler) Start(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)

	// Startup catchup: handle jobs that were due while the service was down.
	if err := s.catchup(ctx); err != nil {
		s.logger.Error("scheduler catchup failed", "error", err)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(time.Duration(s.heartbeatMs) * time.Millisecond)
		defer ticker.Stop()

		s.logger.Info("cron scheduler started", "heartbeatMs", s.heartbeatMs)

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("cron scheduler stopping")
				return
			case <-ticker.C:
				s.tick(ctx)
			}
		}
	}()

	return nil
}

// Stop gracefully stops the scheduler and waits for running jobs to complete.
func (s *Scheduler) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}

	// Wait for all running goroutines with a deadline.
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("cron scheduler stopped gracefully")
	case <-ctx.Done():
		s.logger.Warn("cron scheduler stop timed out, some jobs may still be running")
	}

	return nil
}

// RunningCount returns the number of currently executing jobs.
func (s *Scheduler) RunningCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.running)
}

// catchup handles jobs that were due while the service was down.
// Jobs within the catchup window are triggered immediately; jobs beyond the
// window get their nextRunAtMs recalculated.
func (s *Scheduler) catchup(ctx context.Context) error {
	nowMs := time.Now().UnixMilli()
	cutoffMs := nowMs - s.catchupMs

	// Find all enabled jobs with a past nextRunAtMs.
	dueJobs, err := s.cronRepo.FindDueJobs(ctx, nowMs)
	if err != nil {
		return fmt.Errorf("find due jobs for catchup: %w", err)
	}

	for _, job := range dueJobs {
		nextRunAtMs := job.State().NextRunAtMs()
		if nextRunAtMs == nil {
			continue
		}

		if *nextRunAtMs >= cutoffMs {
			// Within catchup window — trigger immediately.
			s.logger.Info("catchup: triggering overdue job",
				"jobId", job.ID().Value(),
				"name", job.Name(),
				"overdueMs", nowMs-*nextRunAtMs,
			)
			s.launchJob(ctx, job)
		} else {
			// Beyond catchup window — recalculate nextRunAtMs.
			newNext, err := s.schedulerSvc.ComputeNextRun(job.Schedule(), job.CreatedAtMs(), nowMs)
			if err != nil {
				s.logger.Error("catchup: compute next run failed",
					"jobId", job.ID().Value(),
					"error", err,
				)
				continue
			}
			newState := job.State().WithNextRunAtMs(newNext)
			if err := s.cronRepo.UpdateState(ctx, job.ID(), newState); err != nil {
				s.logger.Error("catchup: update state failed",
					"jobId", job.ID().Value(),
					"error", err,
				)
			}
		}
	}

	return nil
}

// tick is called on each heartbeat. It finds due jobs and launches them.
func (s *Scheduler) tick(ctx context.Context) {
	nowMs := time.Now().UnixMilli()

	dueJobs, err := s.cronRepo.FindDueJobs(ctx, nowMs)
	if err != nil {
		s.logger.Error("tick: find due jobs failed", "error", err)
		return
	}

	for _, job := range dueJobs {
		s.launchJob(ctx, job)
	}
}

// launchJob sets the optimistic lock (runningAtMs) and launches async execution.
func (s *Scheduler) launchJob(ctx context.Context, job *cron_job.CronJob) {
	jobID := job.ID().Value()

	// Check if already running in this instance.
	s.mu.Lock()
	if s.running[jobID] {
		s.mu.Unlock()
		return
	}
	s.running[jobID] = true
	s.mu.Unlock()

	// Set runningAtMs as optimistic lock.
	nowMs := time.Now().UnixMilli()
	job.MarkRunning(nowMs)
	if err := s.cronRepo.UpdateState(ctx, job.ID(), job.State()); err != nil {
		s.logger.Warn("failed to set runningAtMs (another instance may have claimed it)",
			"jobId", jobID, "error", err,
		)
		s.mu.Lock()
		delete(s.running, jobID)
		s.mu.Unlock()
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			s.mu.Lock()
			delete(s.running, jobID)
			s.mu.Unlock()
		}()
		s.executeJob(ctx, job)
	}()
}

// executeJob runs the job payload, delivers results, handles alerts, and updates state.
func (s *Scheduler) executeJob(ctx context.Context, job *cron_job.CronJob) {
	jobID := job.ID().Value()
	startMs := time.Now().UnixMilli()

	s.logger.Info("executing cron job", "jobId", jobID, "name", job.Name())

	// Execute the payload.
	result, execErr := s.executor.Execute(ctx, job)
	if execErr != nil {
		result = &ExecuteResult{
			Status: "error",
			Error:  execErr.Error(),
		}
	}

	durationMs := time.Now().UnixMilli() - startMs

	// Deliver results.
	deliveryStatus := "skipped"
	deliveryError := ""
	if result.Status == "ok" {
		ds, dErr := s.delivery.DeliverResult(ctx, job, result)
		deliveryStatus = ds
		if dErr != nil {
			deliveryError = dErr.Error()
		}
	}

	// Handle failure destination for failed jobs.
	if result.Status == "error" && job.Delivery() != nil {
		fd := job.Delivery().FailureDestination()
		if fd != "" && fd != job.Delivery().To() {
			s.logger.Info("sending failure notification to failureDestination",
				"jobId", jobID, "destination", fd,
			)
		}
	}

	// Compute next run time.
	nowMs := time.Now().UnixMilli()
	nextRunAtMs, schedErr := s.schedulerSvc.ComputeNextRun(
		job.Schedule(), job.CreatedAtMs(), nowMs,
	)
	if schedErr != nil {
		s.logger.Error("compute next run failed", "jobId", jobID, "error", schedErr)
	}

	// Update the aggregate state.
	job.MarkCompleted(result.Status, durationMs, nextRunAtMs, result.Error)

	// Update delivery status in state.
	newState := job.State().
		WithLastDeliveryStatus(deliveryStatus).
		WithLastDeliveryError(deliveryError)
	job.UpdateState(newState)

	// Check failure alert.
	s.delivery.SendFailureAlert(ctx, job)

	// Persist updated state.
	if err := s.cronRepo.UpdateState(ctx, job.ID(), job.State()); err != nil {
		s.logger.Error("update job state failed", "jobId", jobID, "error", err)
	}

	// Write run log.
	runLog := &cron_repository.CronRunLog{
		JobID:          job.ID(),
		Ts:             time.Now().UnixMilli(),
		Action:         "finished",
		Status:         result.Status,
		Error:          result.Error,
		Summary:        result.Summary,
		SessionID:      result.SessionID,
		RunAtMs:        startMs,
		DurationMs:     durationMs,
		NextRunAtMs:    nextRunAtMs,
		Model:          result.Model,
		Provider:       result.Provider,
		InputTokens:    result.Usage.InputTokens,
		OutputTokens:   result.Usage.OutputTokens,
		TotalTokens:    result.Usage.TotalTokens,
		DeliveryStatus: deliveryStatus,
		DeliveryError:  deliveryError,
	}
	if err := s.runLogRepo.Append(ctx, runLog); err != nil {
		s.logger.Error("append run log failed", "jobId", jobID, "error", err)
	}

	// Handle deleteAfterRun for one-shot (at) jobs.
	if result.Status == "ok" && job.ShouldDeleteAfterRun() {
		if err := s.cronRepo.Delete(ctx, job.ID()); err != nil {
			s.logger.Error("delete job after run failed", "jobId", jobID, "error", err)
		} else {
			s.logger.Info("deleted job after successful run", "jobId", jobID)
		}
	} else if job.Schedule().Kind() == vo.ScheduleKindAt && nextRunAtMs == nil {
		// At-type job with no next run — disable it.
		disabledState := job.State().WithNextRunAtMs(nil)
		job.UpdateState(disabledState)
		job.Disable()
		if err := s.cronRepo.Update(ctx, job); err != nil {
			s.logger.Error("disable at-type job failed", "jobId", jobID, "error", err)
		}
	}

	s.logger.Info("cron job execution complete",
		"jobId", jobID,
		"status", result.Status,
		"durationMs", durationMs,
	)
}
