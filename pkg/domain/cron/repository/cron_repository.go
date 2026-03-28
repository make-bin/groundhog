package repository

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
)

// CronListFilter is a type-safe filter struct for querying cron jobs.
type CronListFilter struct {
	Enabled       *bool
	AgentId       *string
	SessionTarget *string
	Query         string
	SortBy        string
	SortDir       string
	Offset        int
	Limit         int
}

// CronRepository defines the data access contract for the CronJob aggregate.
type CronRepository interface {
	Create(ctx context.Context, job *cron_job.CronJob) error
	FindByID(ctx context.Context, id vo.CronJobID) (*cron_job.CronJob, error)
	Update(ctx context.Context, job *cron_job.CronJob) error
	Delete(ctx context.Context, id vo.CronJobID) error
	List(ctx context.Context, filter CronListFilter) ([]*cron_job.CronJob, int, error)
	UpdateState(ctx context.Context, id vo.CronJobID, state vo.CronJobState) error
	FindDueJobs(ctx context.Context, nowMs int64) ([]*cron_job.CronJob, error)
}
