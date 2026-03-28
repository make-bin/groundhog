// @AI_GENERATED
package vo

import "fmt"

// PrincipalID is a value object representing the unique identifier of a principal.
// It is immutable after creation.
type PrincipalID struct {
	value string
}

// NewPrincipalID creates a new PrincipalID after validating that the value is non-empty.
func NewPrincipalID(v string) (PrincipalID, error) {
	if v == "" {
		return PrincipalID{}, fmt.Errorf("principal id must not be empty")
	}
	return PrincipalID{value: v}, nil
}

// Value returns the string value of the PrincipalID.
func (id PrincipalID) Value() string { return id.value }

// Equals returns true if id and other represent the same principal identifier.
func (id PrincipalID) Equals(other PrincipalID) bool {
	return id.value == other.value
}

// @AI_GENERATED: end
