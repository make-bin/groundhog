package conversation

import "errors"

var (
	ErrSessionNotFound        = errors.New("session not found")
	ErrSessionAlreadyArchived = errors.New("session already archived")
	ErrModelNotAvailable      = errors.New("model not available")
	ErrToolDenied             = errors.New("tool execution denied")
)
