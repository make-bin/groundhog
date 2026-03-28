package service

import (
	"context"
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/application/assembler"
	"github.com/make-bin/groundhog/pkg/application/dto"
	session_repository "github.com/make-bin/groundhog/pkg/domain/conversation/repository"
	cron_domain "github.com/make-bin/groundhog/pkg/domain/cron"
	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	cron_repository "github.com/make-bin/groundhog/pkg/domain/cron/repository"
	cron_service "github.com/make-bin/groundhog/pkg/domain/cron/service"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// CronAppService defines the application service interface for cron job management.
type CronAppService interface {
	CreateJob(ctx context.Context, req *dto.CreateCronJobRequest) (*dto.CronJobResponse, error)
	UpdateJob(ctx context.Context, id string, req *dto.UpdateCronJobRequest) (*dto.CronJobResponse, error)
	DeleteJob(ctx context.Context, id string) error
	GetJob(ctx context.Context, id string) (*dto.CronJobResponse, error)
	ListJobs(ctx context.Context, req *dto.ListCronJobsRequest) (*dto.CronJobListResponse, error)
	TriggerJob(ctx context.Context, id string, mode string) error
	GetStatus(ctx context.Context) (*dto.CronSchedulerStatusResponse, error)
	GetRunLogs(ctx context.Context, req *dto.GetRunLogsRequest) (*dto.RunLogListResponse, error)
	GetAllRunLogs(ctx context.Context, req *dto.GetRunLogsRequest) (*dto.RunLogListResponse, error)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	// SetScheduler wires the scheduler implementation (called during DI setup to avoid import cycles).
	SetScheduler(scheduler interface {
		Start(context.Context) error
		Stop(context.Context) error
		RunningCount() int
	})
	// SetReaper wires the reaper implementation (called during DI setup to avoid import cycles).
	SetReaper(reaper interface {
		Start(context.Context)
		Stop()
	})
	// Getters for injected dependencies — used by gateway to build Scheduler/Reaper
	// without creating import cycles (infrastructure/cron imports application/service).
	GetAgentAppSvc() AgentAppService
	GetSessionRepo() session_repository.SessionRepository
	GetChannelAppSvc() ChannelAppService
	GetCronRepo() cron_repository.CronRepository
	GetRunLogRepo() cron_repository.CronRunLogRepository
	GetSchedulerSvc() cron_service.CronSchedulerService
}

type cronAppService struct {
	CronRepo      cron_repository.CronRepository       `inject:""`
	RunLogRepo    cron_repository.CronRunLogRepository `inject:""`
	SchedulerSvc  cron_service.CronSchedulerService    `inject:""`
	AgentAppSvc   AgentAppService                      `inject:""`
	SessionRepo   session_repository.SessionRepository `inject:""`
	ChannelAppSvc ChannelAppService                    `inject:""`
	Logger        logger.Logger                        `inject:"logger"`
	scheduler     interface {
		Start(context.Context) error
		Stop(context.Context) error
		RunningCount() int
	}
	reaper interface {
		Start(context.Context)
		Stop()
	}
}

// NewCronAppService creates a new CronAppService. Dependencies are injected via struct tags.
func NewCronAppService() CronAppService {
	return &cronAppService{}
}

// scheduleDTOToVO converts a CronScheduleDTO to a CronSchedule value object.
func scheduleDTOToVO(s dto.CronScheduleDTO) (vo.CronSchedule, error) {
	switch s.Kind {
	case "at":
		return vo.NewCronScheduleAt(s.At)
	case "every":
		return vo.NewCronScheduleEvery(s.EveryMs, s.AnchorMs)
	case "cron":
		return vo.NewCronScheduleCron(s.Expr, s.Tz, s.StaggerMs)
	default:
		return vo.CronSchedule{}, fmt.Errorf("unknown schedule kind: %q", s.Kind)
	}
}

// payloadDTOToVO converts a CronPayloadDTO to a CronPayload value object.
func payloadDTOToVO(p dto.CronPayloadDTO) (vo.CronPayload, error) {
	switch p.Kind {
	case "systemEvent":
		return vo.NewCronPayloadSystemEvent(p.Text)
	case "agentTurn":
		return vo.NewCronPayloadAgentTurn(p.Message, p.Model, p.Thinking, p.TimeoutSeconds, p.LightContext)
	default:
		return vo.CronPayload{}, fmt.Errorf("unknown payload kind: %q", p.Kind)
	}
}

// deliveryDTOToVO converts a CronDeliveryDTO pointer to a CronDelivery value object pointer.
func deliveryDTOToVO(d *dto.CronDeliveryDTO) (*vo.CronDelivery, error) {
	if d == nil {
		return nil, nil
	}
	delivery, err := vo.NewCronDelivery(vo.DeliveryMode(d.Mode), d.Channel, d.To, d.AccountId, d.BestEffort, d.FailureDestination)
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

// failureAlertDTOToVO converts a CronFailureAlertDTO pointer to a CronFailureAlert value object pointer.
func failureAlertDTOToVO(fa *dto.CronFailureAlertDTO) (*vo.CronFailureAlert, error) {
	if fa == nil {
		return nil, nil
	}
	alert, err := vo.NewCronFailureAlert(fa.After, fa.Channel, fa.To, fa.CooldownMs, fa.Mode, fa.AccountId)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// CreateJob creates a new cron job from the request DTO.
func (s *cronAppService) CreateJob(ctx context.Context, req *dto.CreateCronJobRequest) (*dto.CronJobResponse, error) {
	schedule, err := scheduleDTOToVO(req.Schedule)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule: %w", err)
	}

	payload, err := payloadDTOToVO(req.Payload)
	if err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	delivery, err := deliveryDTOToVO(req.Delivery)
	if err != nil {
		return nil, fmt.Errorf("invalid delivery: %w", err)
	}

	failureAlert, err := failureAlertDTOToVO(req.FailureAlert)
	if err != nil {
		return nil, fmt.Errorf("invalid failure_alert: %w", err)
	}

	id := vo.GenerateCronJobID()
	nowMs := time.Now().UnixMilli()
	createdAtMs := nowMs

	// Compute initial nextRunAtMs based on wakeMode
	var computeNowMs int64
	wakeMode := req.WakeMode
	if wakeMode == "" {
		wakeMode = "next-heartbeat"
	}
	if wakeMode == "now" {
		computeNowMs = nowMs - 1 // subtract 1ms so it's immediately due
	} else {
		computeNowMs = nowMs
	}

	nextRunAtMs, err := s.SchedulerSvc.ComputeNextRun(schedule, createdAtMs, computeNowMs)
	if err != nil {
		return nil, fmt.Errorf("compute next run: %w", err)
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	state := vo.NewCronJobState().WithNextRunAtMs(nextRunAtMs)

	job, err := cron_job.NewCronJob(cron_job.CreateCronJobParams{
		ID:             id,
		AgentId:        req.AgentId,
		SessionKey:     req.SessionKey,
		Name:           req.Name,
		Description:    req.Description,
		Enabled:        enabled,
		DeleteAfterRun: req.DeleteAfterRun,
		CreatedAtMs:    createdAtMs,
		UpdatedAtMs:    createdAtMs,
		Schedule:       schedule,
		SessionTarget:  req.SessionTarget,
		WakeMode:       wakeMode,
		Payload:        payload,
		Delivery:       delivery,
		FailureAlert:   failureAlert,
		State:          state,
	})
	if err != nil {
		return nil, err
	}

	if err := s.CronRepo.Create(ctx, job); err != nil {
		return nil, err
	}

	return assembler.ToCronJobResponse(job), nil
}

// UpdateJob updates an existing cron job with the patch from the request DTO.
func (s *cronAppService) UpdateJob(ctx context.Context, id string, req *dto.UpdateCronJobRequest) (*dto.CronJobResponse, error) {
	jobID, err := vo.NewCronJobID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid job id: %w", err)
	}

	job, err := s.CronRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	patch := cron_job.UpdateCronJobParams{
		Name:           req.Name,
		Description:    req.Description,
		Enabled:        req.Enabled,
		DeleteAfterRun: req.DeleteAfterRun,
		SessionTarget:  req.SessionTarget,
		WakeMode:       req.WakeMode,
	}

	// Convert schedule DTO if provided
	if req.Schedule != nil {
		sched, err := scheduleDTOToVO(*req.Schedule)
		if err != nil {
			return nil, fmt.Errorf("invalid schedule: %w", err)
		}
		patch.Schedule = &sched

		// Recompute nextRunAtMs when schedule changes
		nowMs := time.Now().UnixMilli()
		nextRunAtMs, err := s.SchedulerSvc.ComputeNextRun(sched, job.CreatedAtMs(), nowMs)
		if err != nil {
			return nil, fmt.Errorf("compute next run: %w", err)
		}
		newState := job.State().WithNextRunAtMs(nextRunAtMs)
		job.UpdateState(newState)
	}

	// Convert payload DTO if provided
	if req.Payload != nil {
		p, err := payloadDTOToVO(*req.Payload)
		if err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		patch.Payload = &p
	}

	// Convert delivery DTO if provided
	if req.Delivery != nil {
		d, err := deliveryDTOToVO(req.Delivery)
		if err != nil {
			return nil, fmt.Errorf("invalid delivery: %w", err)
		}
		patch.Delivery = d
	}

	// Convert failureAlert DTO if provided
	if req.FailureAlert != nil {
		fa, err := failureAlertDTOToVO(req.FailureAlert)
		if err != nil {
			return nil, fmt.Errorf("invalid failure_alert: %w", err)
		}
		patch.FailureAlert = fa
	}

	if err := job.ApplyPatch(patch); err != nil {
		return nil, err
	}

	if err := s.CronRepo.Update(ctx, job); err != nil {
		return nil, err
	}

	return assembler.ToCronJobResponse(job), nil
}

// DeleteJob soft-deletes a cron job by ID.
func (s *cronAppService) DeleteJob(ctx context.Context, id string) error {
	jobID, err := vo.NewCronJobID(id)
	if err != nil {
		return fmt.Errorf("invalid job id: %w", err)
	}
	return s.CronRepo.Delete(ctx, jobID)
}

// GetJob retrieves a single cron job by ID.
func (s *cronAppService) GetJob(ctx context.Context, id string) (*dto.CronJobResponse, error) {
	jobID, err := vo.NewCronJobID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid job id: %w", err)
	}
	job, err := s.CronRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	return assembler.ToCronJobResponse(job), nil
}

// ListJobs returns a paginated list of cron jobs matching the filter.
func (s *cronAppService) ListJobs(ctx context.Context, req *dto.ListCronJobsRequest) (*dto.CronJobListResponse, error) {
	filter := cron_repository.CronListFilter{
		Query:   req.Query,
		SortBy:  req.SortBy,
		SortDir: req.SortDir,
		Offset:  req.Offset,
		Limit:   req.Limit,
	}

	// Apply enabled filter: if Enabled is explicitly set, use it; otherwise if IncludeDisabled is false, filter to enabled only
	if req.Enabled != nil {
		filter.Enabled = req.Enabled
	} else if !req.IncludeDisabled {
		t := true
		filter.Enabled = &t
	}

	jobs, total, err := s.CronRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return assembler.ToCronJobListResponse(jobs, total, req.Offset, req.Limit), nil
}

// TriggerJob triggers a cron job manually.
// mode "due": only trigger if nextRunAtMs <= nowMs
// mode "force": trigger regardless of schedule
func (s *cronAppService) TriggerJob(ctx context.Context, id string, mode string) error {
	jobID, err := vo.NewCronJobID(id)
	if err != nil {
		return fmt.Errorf("invalid job id: %w", err)
	}

	job, err := s.CronRepo.FindByID(ctx, jobID)
	if err != nil {
		return err
	}

	nowMs := time.Now().UnixMilli()

	if mode == "due" {
		nextRunAtMs := job.State().NextRunAtMs()
		if nextRunAtMs == nil || *nextRunAtMs > nowMs {
			return nil // not due yet
		}
	}

	if job.State().RunningAtMs() != nil {
		return cron_domain.ErrCronJobAlreadyRunning
	}

	job.MarkRunning(nowMs)

	if err := s.CronRepo.UpdateState(ctx, jobID, job.State()); err != nil {
		return err
	}

	return nil
}

// GetStatus returns the overall cron scheduler status.
func (s *cronAppService) GetStatus(ctx context.Context) (*dto.CronSchedulerStatusResponse, error) {
	t := true
	filter := cron_repository.CronListFilter{
		Enabled: &t,
		Limit:   10000, // fetch all enabled jobs
	}

	jobs, _, err := s.CronRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	enabledCount := len(jobs)
	runningCount := 0
	var nextRunAtMs *int64

	for _, job := range jobs {
		state := job.State()
		if state.RunningAtMs() != nil {
			runningCount++
		} else if state.NextRunAtMs() != nil {
			if nextRunAtMs == nil || *state.NextRunAtMs() < *nextRunAtMs {
				v := *state.NextRunAtMs()
				nextRunAtMs = &v
			}
		}
	}

	return &dto.CronSchedulerStatusResponse{
		Running:     s.scheduler != nil,
		EnabledJobs: enabledCount,
		RunningJobs: runningCount,
		NextRunAtMs: nextRunAtMs,
		HeartbeatMs: 10000, // default heartbeat
	}, nil
}

// GetRunLogs returns run logs for a specific job.
func (s *cronAppService) GetRunLogs(ctx context.Context, req *dto.GetRunLogsRequest) (*dto.RunLogListResponse, error) {
	jobID, err := vo.NewCronJobID(req.JobID)
	if err != nil {
		return nil, fmt.Errorf("invalid job id: %w", err)
	}

	filter := cron_repository.RunLogFilter{
		JobID:   &jobID,
		Query:   req.Query,
		SortDir: req.SortDir,
		Offset:  req.Offset,
		Limit:   req.Limit,
	}

	logs, total, err := s.RunLogRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return assembler.ToRunLogListResponse(logs, total, req.Offset, req.Limit), nil
}

// GetAllRunLogs returns run logs across all jobs.
func (s *cronAppService) GetAllRunLogs(ctx context.Context, req *dto.GetRunLogsRequest) (*dto.RunLogListResponse, error) {
	filter := cron_repository.RunLogFilter{
		Query:   req.Query,
		SortDir: req.SortDir,
		Offset:  req.Offset,
		Limit:   req.Limit,
	}

	logs, total, err := s.RunLogRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return assembler.ToRunLogListResponse(logs, total, req.Offset, req.Limit), nil
}

// Start starts the cron scheduler and reaper if they have been wired via SetScheduler/SetReaper.
// The scheduler and reaper are created externally (in the DI wiring layer) to avoid import cycles:
// infrastructure/cron imports application/service, so application/service cannot import infrastructure/cron.
func (s *cronAppService) Start(ctx context.Context) error {
	if s.scheduler == nil || s.reaper == nil {
		s.Logger.Warn("cron scheduler/reaper not configured; wire via SetScheduler/SetReaper before calling Start")
		return nil
	}

	if err := s.scheduler.Start(ctx); err != nil {
		return fmt.Errorf("start scheduler: %w", err)
	}
	s.reaper.Start(ctx)

	s.Logger.Info("cron app service started")
	return nil
}

// Stop gracefully stops the Scheduler and Reaper.
func (s *cronAppService) Stop(ctx context.Context) error {
	if s.reaper != nil {
		s.reaper.Stop()
	}

	if s.scheduler != nil {
		// Give running jobs up to 30 seconds to finish.
		stopCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := s.scheduler.Stop(stopCtx); err != nil {
			s.Logger.Warn("scheduler stop returned error", "error", err)
		}
	}

	s.Logger.Info("cron app service stopped")
	return nil
}

// SetScheduler injects the scheduler implementation (called during DI wiring in gateway).
func (s *cronAppService) SetScheduler(scheduler interface {
	Start(context.Context) error
	Stop(context.Context) error
	RunningCount() int
}) {
	s.scheduler = scheduler
}

// SetReaper injects the reaper implementation (called during DI wiring in gateway).
func (s *cronAppService) SetReaper(reaper interface {
	Start(context.Context)
	Stop()
}) {
	s.reaper = reaper
}

// GetAgentAppSvc returns the injected AgentAppService.
func (s *cronAppService) GetAgentAppSvc() AgentAppService { return s.AgentAppSvc }

// GetSessionRepo returns the injected SessionRepository.
func (s *cronAppService) GetSessionRepo() session_repository.SessionRepository { return s.SessionRepo }

// GetChannelAppSvc returns the injected ChannelAppService.
func (s *cronAppService) GetChannelAppSvc() ChannelAppService { return s.ChannelAppSvc }

// GetCronRepo returns the injected CronRepository.
func (s *cronAppService) GetCronRepo() cron_repository.CronRepository { return s.CronRepo }

// GetRunLogRepo returns the injected CronRunLogRepository.
func (s *cronAppService) GetRunLogRepo() cron_repository.CronRunLogRepository { return s.RunLogRepo }

// GetSchedulerSvc returns the injected CronSchedulerService.
func (s *cronAppService) GetSchedulerSvc() cron_service.CronSchedulerService { return s.SchedulerSvc }
