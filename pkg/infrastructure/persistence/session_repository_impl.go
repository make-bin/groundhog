// @AI_GENERATED
package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	conversation "github.com/make-bin/groundhog/pkg/domain/conversation"
	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
	"github.com/make-bin/groundhog/pkg/domain/conversation/repository"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/mapper"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

type sessionRepositoryImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
}

// NewSessionRepository creates a new SessionRepository implementation.
func NewSessionRepository() repository.SessionRepository {
	return &sessionRepositoryImpl{}
}

// Create persists a new AgentSession aggregate.
func (r *sessionRepositoryImpl) Create(ctx context.Context, session *agent_session.AgentSession) error {
	sessionPO, err := mapper.DomainToSessionPO(session)
	if err != nil {
		return err
	}
	return r.DataStore.DB().WithContext(ctx).Create(sessionPO).Error
}

// FindByID retrieves a complete AgentSession aggregate by its SessionID.
func (r *sessionRepositoryImpl) FindByID(ctx context.Context, id vo.SessionID) (*agent_session.AgentSession, error) {
	var sessionPO po.SessionPO
	result := r.DataStore.DB().WithContext(ctx).
		Preload("Turns").
		Where("session_id = ?", id.Value()).
		First(&sessionPO)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, conversation.ErrSessionNotFound
		}
		return nil, result.Error
	}
	return mapper.SessionPOToDomain(&sessionPO)
}

// Update persists changes to an existing AgentSession aggregate.
func (r *sessionRepositoryImpl) Update(ctx context.Context, session *agent_session.AgentSession) error {
	sessionPO, err := mapper.DomainToSessionPO(session)
	if err != nil {
		return err
	}

	// Fetch the existing row to get the auto-increment ID
	var existing po.SessionPO
	if err := r.DataStore.DB().WithContext(ctx).
		Select("id").
		Where("session_id = ?", session.ID().Value()).
		First(&existing).Error; err != nil {
		return err
	}

	// Update scalar fields
	if err := r.DataStore.DB().WithContext(ctx).
		Model(&po.SessionPO{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"active_model": sessionPO.ActiveModel,
			"state":        sessionPO.State,
			"metadata":     sessionPO.Metadata,
		}).Error; err != nil {
		return err
	}

	// Replace turns: delete old, insert new
	if err := r.DataStore.DB().WithContext(ctx).
		Where("session_po_id = ?", existing.ID).
		Delete(&po.TurnPO{}).Error; err != nil {
		return err
	}
	if len(sessionPO.Turns) > 0 {
		for i := range sessionPO.Turns {
			sessionPO.Turns[i].SessionPOID = existing.ID
		}
		if err := r.DataStore.DB().WithContext(ctx).Create(&sessionPO.Turns).Error; err != nil {
			return err
		}
	}

	return nil
}

// Delete removes an AgentSession aggregate by its SessionID.
func (r *sessionRepositoryImpl) Delete(ctx context.Context, id vo.SessionID) error {
	return r.DataStore.DB().WithContext(ctx).
		Where("session_id = ?", id.Value()).
		Delete(&po.SessionPO{}).Error
}

// List retrieves AgentSession aggregates matching the filter with pagination.
func (r *sessionRepositoryImpl) List(ctx context.Context, filter repository.SessionFilter, offset, limit int) ([]*agent_session.AgentSession, int, error) {
	db := r.DataStore.DB().WithContext(ctx).Model(&po.SessionPO{})

	if filter.UserID != nil {
		db = db.Where("user_id = ?", *filter.UserID)
	}
	if filter.AgentID != nil {
		db = db.Where("agent_id = ?", filter.AgentID.Value())
	}
	if filter.State != nil {
		db = db.Where("state = ?", int(*filter.State))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var sessionPOs []po.SessionPO
	if err := db.Preload("Turns").Offset(offset).Limit(limit).Find(&sessionPOs).Error; err != nil {
		return nil, 0, err
	}

	sessions := make([]*agent_session.AgentSession, 0, len(sessionPOs))
	for i := range sessionPOs {
		s, err := mapper.SessionPOToDomain(&sessionPOs[i])
		if err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, s)
	}
	return sessions, int(total), nil
}

// Archive transitions an AgentSession to the archived state.
func (r *sessionRepositoryImpl) Archive(ctx context.Context, id vo.SessionID) error {
	return r.DataStore.DB().WithContext(ctx).
		Model(&po.SessionPO{}).
		Where("session_id = ?", id.Value()).
		Update("state", int(vo.SessionStateArchived)).Error
}

// @AI_GENERATED: end
