package vo

import "errors"

// MediaRef is a reference to a media asset by its ID.
type MediaRef struct {
	assetID string
}

var ErrEmptyMediaRef = errors.New("media ref asset ID cannot be empty")

func NewMediaRef(assetID string) (MediaRef, error) {
	if assetID == "" {
		return MediaRef{}, ErrEmptyMediaRef
	}
	return MediaRef{assetID: assetID}, nil
}

func (r MediaRef) AssetID() string            { return r.assetID }
func (r MediaRef) Equals(other MediaRef) bool { return r.assetID == other.assetID }
