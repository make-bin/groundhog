// @AI_GENERATED
package vo

import "fmt"

// AccountID represents a unique identifier for an account (user) on a channel.
// It is immutable after creation.
type AccountID struct {
	value string
}

// NewAccountID creates a new AccountID after validating that the value is non-empty.
func NewAccountID(value string) (AccountID, error) {
	if value == "" {
		return AccountID{}, fmt.Errorf("account ID must not be empty")
	}
	return AccountID{value: value}, nil
}

// Value returns the account ID string.
func (a AccountID) Value() string { return a.value }

// Equals returns true if a and other represent the same account ID.
func (a AccountID) Equals(other AccountID) bool {
	return a.value == other.value
}

// @AI_GENERATED: end
