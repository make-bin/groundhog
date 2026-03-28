// @AI_GENERATED
package vo

import "fmt"

// SessionID represents a unique identifier for an agent session.
// It is immutable after creation.
type SessionID struct {
	value string
}

// NewSessionID creates a new SessionID after validating that the value is non-empty.
func NewSessionID(value string) (SessionID, error) {
	if value == "" {
		return SessionID{}, fmt.Errorf("session ID must not be empty")
	}
	return SessionID{value: value}, nil
}

// Value returns the session ID string.
func (s SessionID) Value() string { return s.value }

// Equals returns true if s and other represent the same session ID.
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}

// @AI_GENERATED: end
