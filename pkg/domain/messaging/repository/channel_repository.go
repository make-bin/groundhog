package repository

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/messaging/entity"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// ChannelRepository defines the data access contract for Channel entities (managed as aggregate).
type ChannelRepository interface {
	FindAll(ctx context.Context) ([]*entity.Channel, error)
	FindByID(ctx context.Context, id vo.ChannelID) (*entity.Channel, error)
	Save(ctx context.Context, channel *entity.Channel) error
	Delete(ctx context.Context, id vo.ChannelID) error
}
