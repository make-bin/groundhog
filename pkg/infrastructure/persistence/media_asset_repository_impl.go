package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/make-bin/groundhog/pkg/domain/media/aggregate/media_asset"
	"github.com/make-bin/groundhog/pkg/domain/media/repository"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/mapper"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

type mediaAssetRepositoryImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
}

// NewMediaAssetRepository creates a new MediaRepository implementation.
func NewMediaAssetRepository() repository.MediaRepository {
	return &mediaAssetRepositoryImpl{}
}

func (r *mediaAssetRepositoryImpl) Save(ctx context.Context, asset *media_asset.MediaAsset) error {
	p := mapper.DomainToMediaAssetPO(asset)
	return r.DataStore.DB().WithContext(ctx).Save(p).Error
}

func (r *mediaAssetRepositoryImpl) FindByID(ctx context.Context, id string) (*media_asset.MediaAsset, error) {
	var p po.MediaAssetPO
	result := r.DataStore.DB().WithContext(ctx).Where("asset_id = ?", id).First(&p)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, media_asset.ErrMediaAssetNotFound
		}
		return nil, result.Error
	}
	return mapper.MediaAssetPOToDomain(&p)
}

func (r *mediaAssetRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.DataStore.DB().WithContext(ctx).
		Where("asset_id = ?", id).
		Delete(&po.MediaAssetPO{}).Error
}
