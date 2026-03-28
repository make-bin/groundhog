// @AI_GENERATED
package vo

import "fmt"

// ChannelID represents a unique identifier for a messaging channel.
// It is immutable after creation.
type ChannelID struct {
	value string
}

// NewChannelID creates a new ChannelID after validating that the value is non-empty.
func NewChannelID(value string) (ChannelID, error) {
	if value == "" {
		return ChannelID{}, fmt.Errorf("channel ID must not be empty")
	}
	return ChannelID{value: value}, nil
}

// Value returns the channel ID string.
func (c ChannelID) Value() string { return c.value }

// Equals returns true if c and other represent the same channel ID.
func (c ChannelID) Equals(other ChannelID) bool {
	return c.value == other.value
}

// @AI_GENERATED: end
