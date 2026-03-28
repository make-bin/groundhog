// @AI_GENERATED
package vo

import "fmt"

// AgentID represents a unique identifier for an agent.
// It is immutable after creation.
type AgentID struct {
	value string
}

// NewAgentID creates a new AgentID after validating that the value is non-empty.
func NewAgentID(value string) (AgentID, error) {
	if value == "" {
		return AgentID{}, fmt.Errorf("agent ID must not be empty")
	}
	return AgentID{value: value}, nil
}

// Value returns the agent ID string.
func (a AgentID) Value() string { return a.value }

// Equals returns true if a and other represent the same agent ID.
func (a AgentID) Equals(other AgentID) bool {
	return a.value == other.value
}

// @AI_GENERATED: end
