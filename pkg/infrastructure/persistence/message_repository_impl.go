package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	messaging "github.com/make-bin/groundhog/pkg/domain/messaging"
	"github.com/make-bin/groundhog/pkg/domain/messaging/aggregate/inbound_message"
	"github.com/make-bin/groundhog/pkg/domain/messaging/repository"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/mapper"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

type messageRepositoryImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
}

// NewMessageRepository creates a new MessageRepository implementation.
func NewMessageRepository() repository.MessageRepository {
	return &messageRepositoryImpl{}
}

// Create persists a new InboundMessage aggregate.
func (r *messageRepositoryImpl) Create(ctx context.Context, msg *inbound_message.InboundMessage) error {
	msgPO, err := mapper.DomainToMessagePO(msg)
	if err != nil {
		return err
	}
	return r.DataStore.DB().WithContext(ctx).Create(msgPO).Error
}

// FindByID retrieves an InboundMessage aggregate by its MessageID.
func (r *messageRepositoryImpl) FindByID(ctx context.Context, id vo.MessageID) (*inbound_message.InboundMessage, error) {
	var msgPO po.MessagePO
	result := r.DataStore.DB().WithContext(ctx).
		Where("message_id = ?", id.Value()).
		First(&msgPO)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, messaging.ErrMessageNotFound
		}
		return nil, result.Error
	}
	return mapper.MessagePOToDomain(&msgPO)
}

// Update persists changes to an existing InboundMessage aggregate.
func (r *messageRepositoryImpl) Update(ctx context.Context, msg *inbound_message.InboundMessage) error {
	msgPO, err := mapper.DomainToMessagePO(msg)
	if err != nil {
		return err
	}
	return r.DataStore.DB().WithContext(ctx).
		Where("message_id = ?", msgPO.MessageID).
		Save(msgPO).Error
}

// List retrieves InboundMessage aggregates matching the filter with pagination.
func (r *messageRepositoryImpl) List(ctx context.Context, filter repository.MessageFilter, offset, limit int) ([]*inbound_message.InboundMessage, int, error) {
	db := r.DataStore.DB().WithContext(ctx).Model(&po.MessagePO{})

	if filter.ChannelID != nil {
		db = db.Where("channel_id = ?", filter.ChannelID.Value())
	}
	if filter.AccountID != nil {
		db = db.Where("account_id = ?", filter.AccountID.Value())
	}
	if filter.Status != nil {
		db = db.Where("status = ?", int(*filter.Status))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var msgPOs []po.MessagePO
	if err := db.Offset(offset).Limit(limit).Find(&msgPOs).Error; err != nil {
		return nil, 0, err
	}

	msgs := make([]*inbound_message.InboundMessage, 0, len(msgPOs))
	for i := range msgPOs {
		m, err := mapper.MessagePOToDomain(&msgPOs[i])
		if err != nil {
			return nil, 0, err
		}
		msgs = append(msgs, m)
	}
	return msgs, int(total), nil
}
