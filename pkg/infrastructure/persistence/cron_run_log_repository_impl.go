package persistence

import (
	"context"
	"fmt"

	cron_repository "github.com/make-bin/groundhog/pkg/domain/cron/repository"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/mapper"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

type cronRunLogRepositoryImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
}

// NewCronRunLogRepositoryImpl creates a new CronRunLogRepository implementation.
func NewCronRunLogRepositoryImpl() cron_repository.CronRunLogRepository {
	return &cronRunLogRepositoryImpl{}
}

// Append persists a new CronRunLog entry.
func (r *cronRunLogRepositoryImpl) Append(ctx context.Context, log *cron_repository.CronRunLog) error {
	p := mapper.DomainToCronRunLogPO(log)
	return r.DataStore.DB().WithContext(ctx).Create(p).Error
}

// List retrieves CronRunLog entries matching the filter with pagination.
func (r *cronRunLogRepositoryImpl) List(ctx context.Context, filter cron_repository.RunLogFilter) ([]*cron_repository.CronRunLog, int, error) {
	db := r.DataStore.DB().WithContext(ctx).Model(&po.CronRunLogPO{})

	if filter.JobID != nil {
		db = db.Where("job_id = ?", filter.JobID.Value())
	}
	if filter.Status != nil {
		db = db.Where("status = ?", *filter.Status)
	}
	if filter.DeliveryStatus != nil {
		db = db.Where("delivery_status = ?", *filter.DeliveryStatus)
	}
	if filter.Query != "" {
		like := "%" + filter.Query + "%"
		db = db.Where("summary LIKE ? OR error LIKE ?", like, like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortDir := filter.SortDir
	if sortDir == "" {
		sortDir = "DESC"
	}
	db = db.Order(fmt.Sprintf("run_at_ms %s", sortDir))

	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}

	var pos []po.CronRunLogPO
	if err := db.Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	logs := make([]*cron_repository.CronRunLog, 0, len(pos))
	for i := range pos {
		l, err := mapper.CronRunLogPOToDomain(&pos[i])
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, int(total), nil
}
