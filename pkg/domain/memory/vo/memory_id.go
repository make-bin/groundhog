package vo

import "errors"

// MemoryID represents a unique identifier for a memory entry.
// It is immutable after creation.
type MemoryID struct {
	value string
}

// NewMemoryID creates a new MemoryID after validating that the value is non-empty.
func NewMemoryID(v string) (MemoryID, error) {
	if v == "" {
		return MemoryID{}, errors.New("memory id must not be empty")
	}
	return MemoryID{value: v}, nil
}

// Value returns the memory ID string.
func (id MemoryID) Value() string {
	return id.value
}
