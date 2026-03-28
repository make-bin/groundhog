package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	messaging "github.com/make-bin/groundhog/pkg/domain/messaging"
	"github.com/make-bin/groundhog/pkg/domain/messaging/entity"
	"github.com/make-bin/groundhog/pkg/domain/messaging/repository"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/mapper"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

type channelRepositoryImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
}

// NewChannelRepository creates a new ChannelRepository implementation.
func NewChannelRepository() repository.ChannelRepository {
	return &channelRepositoryImpl{}
}

// FindAll retrieves all Channel entities.
func (r *channelRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Channel, error) {
	var channelPOs []po.ChannelPO
	if err := r.DataStore.DB().WithContext(ctx).Find(&channelPOs).Error; err != nil {
		return nil, err
	}
	channels := make([]*entity.Channel, 0, len(channelPOs))
	for i := range channelPOs {
		ch, err := mapper.ChannelPOToDomain(&channelPOs[i])
		if err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}
	return channels, nil
}

// FindByID retrieves a Channel entity by its ChannelID.
func (r *channelRepositoryImpl) FindByID(ctx context.Context, id vo.ChannelID) (*entity.Channel, error) {
	var channelPO po.ChannelPO
	result := r.DataStore.DB().WithContext(ctx).
		Where("channel_id = ?", id.Value()).
		First(&channelPO)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, messaging.ErrChannelNotFound
		}
		return nil, result.Error
	}
	return mapper.ChannelPOToDomain(&channelPO)
}

// Save persists a Channel entity (upsert by channel_id).
func (r *channelRepositoryImpl) Save(ctx context.Context, channel *entity.Channel) error {
	channelPO, err := mapper.DomainToChannelPO(channel)
	if err != nil {
		return err
	}
	return r.DataStore.DB().WithContext(ctx).
		Where("channel_id = ?", channelPO.ChannelID).
		Save(channelPO).Error
}

// Delete removes a Channel entity by its ChannelID.
func (r *channelRepositoryImpl) Delete(ctx context.Context, id vo.ChannelID) error {
	return r.DataStore.DB().WithContext(ctx).
		Where("channel_id = ?", id.Value()).
		Delete(&po.ChannelPO{}).Error
}
