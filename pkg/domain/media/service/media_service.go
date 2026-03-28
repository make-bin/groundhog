package service

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/media/aggregate/media_asset"
	"github.com/make-bin/groundhog/pkg/domain/media/vo"
)

// MediaService defines the domain service for media operations.
type MediaService interface {
	Fetch(ctx context.Context, url string) (*media_asset.MediaAsset, error)
	Detect(data []byte) (vo.MimeType, error)
	Validate(asset *media_asset.MediaAsset) error
}
