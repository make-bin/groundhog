package repository

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/media/aggregate/media_asset"
)

// MediaRepository defines the data access contract for MediaAsset aggregate.
type MediaRepository interface {
	Save(ctx context.Context, asset *media_asset.MediaAsset) error
	FindByID(ctx context.Context, id string) (*media_asset.MediaAsset, error)
	Delete(ctx context.Context, id string) error
}
