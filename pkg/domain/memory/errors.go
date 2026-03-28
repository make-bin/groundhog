package memory

import "errors"

var (
	ErrMemoryNotFound     = errors.New("memory not found")
	ErrMemoryAccessDenied = errors.New("memory access denied")
	ErrEmptyContent       = errors.New("memory content must not be empty")
	ErrMissingUserID      = errors.New("userID is required in context")
)
