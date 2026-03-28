package mapper

import (
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/media/aggregate/media_asset"
	"github.com/make-bin/groundhog/pkg/domain/media/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

// DomainToMediaAssetPO converts a MediaAsset to a MediaAssetPO.
func DomainToMediaAssetPO(asset *media_asset.MediaAsset) *po.MediaAssetPO {
	return &po.MediaAssetPO{
		AssetID:   asset.ID(),
		MimeType:  asset.MimeType().Value(),
		Size:      asset.Size(),
		Source:    int(asset.Source()),
		URL:       asset.URL(),
		LocalPath: asset.LocalPath(),
		Validated: asset.Validated(),
	}
}

// MediaAssetPOToDomain converts a MediaAssetPO to a MediaAsset.
func MediaAssetPOToDomain(p *po.MediaAssetPO) (*media_asset.MediaAsset, error) {
	mimeType, err := vo.NewMimeType(p.MimeType)
	if err != nil {
		return nil, fmt.Errorf("reconstruct mime_type: %w", err)
	}
	return media_asset.ReconstructMediaAsset(
		p.AssetID,
		mimeType,
		p.Size,
		media_asset.MediaSource(p.Source),
		p.URL,
		p.LocalPath,
		p.Validated,
		p.CreatedAt,
		p.UpdatedAt,
	), nil
}
