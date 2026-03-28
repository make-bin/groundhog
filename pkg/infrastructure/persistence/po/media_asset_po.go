package po

import (
	"time"

	"gorm.io/gorm"
)

// MediaAssetPO is the persistence object for MediaAsset aggregate.
type MediaAssetPO struct {
	ID        uint   `gorm:"primaryKey"`
	AssetID   string `gorm:"uniqueIndex;type:varchar(100);not null"`
	MimeType  string `gorm:"type:varchar(100)"`
	Size      int64
	Source    int
	URL       string `gorm:"type:text"`
	LocalPath string `gorm:"type:text"`
	Validated bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (MediaAssetPO) TableName() string { return "media_assets" }
