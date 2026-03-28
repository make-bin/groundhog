package repository

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
)

// RunLogFilter is a type-safe filter struct for querying cron run logs.
type RunLogFilter struct {
	JobID          *vo.CronJobID
	Status         *string
	DeliveryStatus *string
	Query          string
	SortDir        string
	Offset         int
	Limit          int
}

// CronRunLog is the run log entity (not an aggregate root).
type CronRunLog struct {
	ID             int64
	JobID          vo.CronJobID
	Ts             int64
	Action         string
	Status         string // ok/error
	Error          string
	Summary        string
	SessionID      string
	RunAtMs        int64
	DurationMs     int64
	NextRunAtMs    *int64
	Model          string
	Provider       string
	InputTokens    int
	OutputTokens   int
	TotalTokens    int
	DeliveryStatus string
	DeliveryError  string
}

// CronRunLogRepository defines the data access contract for cron run logs.
type CronRunLogRepository interface {
	Append(ctx context.Context, log *CronRunLog) error
	List(ctx context.Context, filter RunLogFilter) ([]*CronRunLog, int, error)
}
