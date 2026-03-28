package vo

import (
	"errors"
	"strings"
)

// MimeType is a validated MIME type value object.
type MimeType struct {
	value string
}

var ErrInvalidMimeType = errors.New("invalid MIME type format")

// NewMimeType creates a MimeType, validating format (must contain "/").
func NewMimeType(value string) (MimeType, error) {
	if !strings.Contains(value, "/") {
		return MimeType{}, ErrInvalidMimeType
	}
	return MimeType{value: strings.ToLower(strings.TrimSpace(value))}, nil
}

func (m MimeType) Value() string              { return m.value }
func (m MimeType) Equals(other MimeType) bool { return m.value == other.value }
func (m MimeType) String() string             { return m.value }
