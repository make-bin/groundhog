package vo

import "fmt"

// Capability represents a capability string that a plugin exposes.
// It is immutable after creation.
type Capability struct {
	value string
}

// NewCapability creates a new Capability after validating that the value is non-empty.
func NewCapability(value string) (Capability, error) {
	if value == "" {
		return Capability{}, fmt.Errorf("capability must not be empty")
	}
	return Capability{value: value}, nil
}

// Value returns the capability string.
func (c Capability) Value() string { return c.value }

// Equals returns true if c and other represent the same capability.
func (c Capability) Equals(other Capability) bool {
	return c.value == other.value
}
