// @AI_GENERATED
package service

import (
	"context"
	"time"

	identitysvc "github.com/make-bin/groundhog/pkg/domain/identity/service"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

type auditServiceImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
	Logger    logger.Logger       `inject:"logger"`
}

// NewAuditService creates a new AuditService implementation.
func NewAuditService() identitysvc.AuditService {
	return &auditServiceImpl{}
}

func (s *auditServiceImpl) Record(ctx context.Context, action, principalID, resourceType, resourceID, details, sourceIP string) error {
	entry := &po.AuditLogPO{
		Action:       action,
		PrincipalID:  principalID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
		SourceIP:     sourceIP,
		CreatedAt:    time.Now(),
	}
	if err := s.DataStore.DB().WithContext(ctx).Create(entry).Error; err != nil {
		s.Logger.Error("failed to record audit log", "action", action, "error", err)
		return err
	}
	return nil
}

func (s *auditServiceImpl) Query(ctx context.Context, filter identitysvc.AuditFilter) ([]*identitysvc.AuditLog, int, error) {
	db := s.DataStore.DB().WithContext(ctx).Model(&po.AuditLogPO{})

	if filter.Action != nil {
		db = db.Where("action = ?", *filter.Action)
	}
	if filter.PrincipalID != nil {
		db = db.Where("principal_id = ?", *filter.PrincipalID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var entries []po.AuditLogPO
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	logs := make([]*identitysvc.AuditLog, 0, len(entries))
	for i := range entries {
		e := &entries[i]
		logs = append(logs, &identitysvc.AuditLog{
			ID:           e.ID,
			Action:       e.Action,
			PrincipalID:  e.PrincipalID,
			ResourceType: e.ResourceType,
			ResourceID:   e.ResourceID,
			Details:      e.Details,
			SourceIP:     e.SourceIP,
			CreatedAt:    e.CreatedAt,
		})
	}
	return logs, int(total), nil
}

// @AI_GENERATED: end
