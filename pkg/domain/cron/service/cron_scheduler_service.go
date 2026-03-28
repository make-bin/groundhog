package service

import "github.com/make-bin/groundhog/pkg/domain/cron/vo"

// CronSchedulerService is the domain service for computing cron schedule times.
type CronSchedulerService interface {
	// ComputeNextRun computes the next run time in ms. Returns nil if no next run (e.g., at-type expired).
	ComputeNextRun(schedule vo.CronSchedule, createdAtMs, nowMs int64) (*int64, error)
	// IsValidSchedule validates the schedule.
	IsValidSchedule(schedule vo.CronSchedule) error
}

// NewCronSchedulerService creates a new CronSchedulerService.
func NewCronSchedulerService() CronSchedulerService {
	return &cronSchedulerServiceImpl{}
}
