package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	cron_domain "github.com/make-bin/groundhog/pkg/domain/cron"
	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	cron_repository "github.com/make-bin/groundhog/pkg/domain/cron/repository"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/mapper"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

type cronJobRepositoryImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
}

// NewCronJobRepositoryImpl creates a new CronRepository implementation.
func NewCronJobRepositoryImpl() cron_repository.CronRepository {
	return &cronJobRepositoryImpl{}
}

// Create persists a new CronJob aggregate.
func (r *cronJobRepositoryImpl) Create(ctx context.Context, job *cron_job.CronJob) error {
	p, err := mapper.DomainToCronJobPO(job)
	if err != nil {
		return fmt.Errorf("map domain to po: %w", err)
	}
	if err := r.DataStore.DB().WithContext(ctx).Create(p).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return cron_domain.ErrCronJobAlreadyExists
		}
		return err
	}
	return nil
}

// FindByID retrieves a CronJob by its job_id.
func (r *cronJobRepositoryImpl) FindByID(ctx context.Context, id vo.CronJobID) (*cron_job.CronJob, error) {
	var p po.CronJobPO
	if err := r.DataStore.DB().WithContext(ctx).
		Where("job_id = ?", id.Value()).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cron_domain.ErrCronJobNotFound
		}
		return nil, err
	}
	return mapper.CronJobPOToDomain(&p)
}

// Update persists changes to an existing CronJob aggregate.
func (r *cronJobRepositoryImpl) Update(ctx context.Context, job *cron_job.CronJob) error {
	p, err := mapper.DomainToCronJobPO(job)
	if err != nil {
		return fmt.Errorf("map domain to po: %w", err)
	}

	var existing po.CronJobPO
	if err := r.DataStore.DB().WithContext(ctx).
		Where("job_id = ?", job.ID().Value()).
		First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cron_domain.ErrCronJobNotFound
		}
		return err
	}
	p.ID = existing.ID
	p.CreatedAt = existing.CreatedAt // preserve original created_at

	return r.DataStore.DB().WithContext(ctx).Save(p).Error
}

// Delete soft-deletes a CronJob by its job_id.
func (r *cronJobRepositoryImpl) Delete(ctx context.Context, id vo.CronJobID) error {
	result := r.DataStore.DB().WithContext(ctx).
		Where("job_id = ?", id.Value()).
		Delete(&po.CronJobPO{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return cron_domain.ErrCronJobNotFound
	}
	return nil
}

// List retrieves CronJob aggregates matching the filter with pagination.
func (r *cronJobRepositoryImpl) List(ctx context.Context, filter cron_repository.CronListFilter) ([]*cron_job.CronJob, int, error) {
	db := r.DataStore.DB().WithContext(ctx).Model(&po.CronJobPO{})

	if filter.Enabled != nil {
		db = db.Where("enabled = ?", *filter.Enabled)
	}
	if filter.AgentId != nil {
		db = db.Where("agent_id = ?", *filter.AgentId)
	}
	if filter.SessionTarget != nil {
		db = db.Where("session_target = ?", *filter.SessionTarget)
	}
	if filter.Query != "" {
		like := "%" + filter.Query + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", like, like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "updated_at"
	}
	sortDir := filter.SortDir
	if sortDir == "" {
		sortDir = "DESC"
	}
	db = db.Order(fmt.Sprintf("%s %s", sortBy, sortDir))

	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}

	var pos []po.CronJobPO
	if err := db.Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	jobs := make([]*cron_job.CronJob, 0, len(pos))
	for i := range pos {
		j, err := mapper.CronJobPOToDomain(&pos[i])
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, j)
	}
	return jobs, int(total), nil
}

// UpdateState updates only the state JSONB field for a CronJob, without overwriting config fields.
func (r *cronJobRepositoryImpl) UpdateState(ctx context.Context, id vo.CronJobID, state vo.CronJobState) error {
	type stateJSON struct {
		NextRunAtMs          *int64 `json:"next_run_at_ms"`
		RunningAtMs          *int64 `json:"running_at_ms"`
		LastRunAtMs          *int64 `json:"last_run_at_ms"`
		LastRunStatus        string `json:"last_run_status"`
		LastError            string `json:"last_error"`
		LastDurationMs       int64  `json:"last_duration_ms"`
		ConsecutiveErrors    int    `json:"consecutive_errors"`
		LastFailureAlertAtMs *int64 `json:"last_failure_alert_at_ms"`
		ScheduleErrorCount   int    `json:"schedule_error_count"`
		LastDeliveryStatus   string `json:"last_delivery_status"`
		LastDeliveryError    string `json:"last_delivery_error"`
	}

	stateBytes, err := json.Marshal(stateJSON{
		NextRunAtMs:          state.NextRunAtMs(),
		RunningAtMs:          state.RunningAtMs(),
		LastRunAtMs:          state.LastRunAtMs(),
		LastRunStatus:        state.LastRunStatus(),
		LastError:            state.LastError(),
		LastDurationMs:       state.LastDurationMs(),
		ConsecutiveErrors:    state.ConsecutiveErrors(),
		LastFailureAlertAtMs: state.LastFailureAlertAtMs(),
		ScheduleErrorCount:   state.ScheduleErrorCount(),
		LastDeliveryStatus:   state.LastDeliveryStatus(),
		LastDeliveryError:    state.LastDeliveryError(),
	})
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}

	result := r.DataStore.DB().WithContext(ctx).
		Model(&po.CronJobPO{}).
		Where("job_id = ?", id.Value()).
		Update("state", gorm.Expr("?::jsonb", string(stateBytes)))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return cron_domain.ErrCronJobNotFound
	}
	return nil
}

// FindDueJobs returns all enabled jobs whose next_run_at_ms is due and are not currently running.
func (r *cronJobRepositoryImpl) FindDueJobs(ctx context.Context, nowMs int64) ([]*cron_job.CronJob, error) {
	var pos []po.CronJobPO
	if err := r.DataStore.DB().WithContext(ctx).
		Where("enabled = ? AND (state->>'running_at_ms') IS NULL AND (state->>'next_run_at_ms') IS NOT NULL AND CAST(state->>'next_run_at_ms' AS BIGINT) <= ?",
			true, nowMs).
		Find(&pos).Error; err != nil {
		return nil, err
	}

	jobs := make([]*cron_job.CronJob, 0, len(pos))
	for i := range pos {
		j, err := mapper.CronJobPOToDomain(&pos[i])
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}
