package vo

import "fmt"

// PluginID represents a unique identifier for a plugin instance.
// It is immutable after creation.
type PluginID struct {
	value string
}

// NewPluginID creates a new PluginID after validating that the value is non-empty.
func NewPluginID(value string) (PluginID, error) {
	if value == "" {
		return PluginID{}, fmt.Errorf("plugin ID must not be empty")
	}
	return PluginID{value: value}, nil
}

// Value returns the plugin ID string.
func (p PluginID) Value() string { return p.value }

// Equals returns true if p and other represent the same plugin ID.
func (p PluginID) Equals(other PluginID) bool {
	return p.value == other.value
}
