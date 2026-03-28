// @AI_GENERATED
package vo

import "fmt"

// MessageID represents a unique identifier for a message.
// It is immutable after creation.
type MessageID struct {
	value string
}

// NewMessageID creates a new MessageID after validating that the value is non-empty.
func NewMessageID(value string) (MessageID, error) {
	if value == "" {
		return MessageID{}, fmt.Errorf("message ID must not be empty")
	}
	return MessageID{value: value}, nil
}

// Value returns the message ID string.
func (m MessageID) Value() string { return m.value }

// Equals returns true if m and other represent the same message ID.
func (m MessageID) Equals(other MessageID) bool {
	return m.value == other.value
}

// @AI_GENERATED: end
