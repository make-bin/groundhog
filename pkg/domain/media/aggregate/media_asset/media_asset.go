package media_asset

import (
	"errors"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/media/vo"
)

// MediaSource indicates where the media came from.
type MediaSource int

const (
	MediaSourceURL    MediaSource = 0
	MediaSourceLocal  MediaSource = 1
	MediaSourceUpload MediaSource = 2
)

// MediaAsset is the aggregate root for a media resource.
type MediaAsset struct {
	id        string
	mimeType  vo.MimeType
	size      int64
	source    MediaSource
	url       string
	localPath string
	validated bool
	createdAt time.Time
	updatedAt time.Time
}

var ErrMediaAssetNotFound = errors.New("media asset not found")

// NewMediaAsset creates a new MediaAsset.
func NewMediaAsset(id string, mimeType vo.MimeType, size int64, source MediaSource, url, localPath string) *MediaAsset {
	now := time.Now()
	return &MediaAsset{
		id:        id,
		mimeType:  mimeType,
		size:      size,
		source:    source,
		url:       url,
		localPath: localPath,
		validated: false,
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructMediaAsset reconstructs a MediaAsset from persistence.
func ReconstructMediaAsset(id string, mimeType vo.MimeType, size int64, source MediaSource, url, localPath string, validated bool, createdAt, updatedAt time.Time) *MediaAsset {
	return &MediaAsset{
		id:        id,
		mimeType:  mimeType,
		size:      size,
		source:    source,
		url:       url,
		localPath: localPath,
		validated: validated,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// Validate marks the asset as validated.
func (a *MediaAsset) Validate() {
	a.validated = true
	a.updatedAt = time.Now()
}

// Getters
func (a *MediaAsset) ID() string            { return a.id }
func (a *MediaAsset) MimeType() vo.MimeType { return a.mimeType }
func (a *MediaAsset) Size() int64           { return a.size }
func (a *MediaAsset) Source() MediaSource   { return a.source }
func (a *MediaAsset) URL() string           { return a.url }
func (a *MediaAsset) LocalPath() string     { return a.localPath }
func (a *MediaAsset) Validated() bool       { return a.validated }
func (a *MediaAsset) CreatedAt() time.Time  { return a.createdAt }
func (a *MediaAsset) UpdatedAt() time.Time  { return a.updatedAt }
