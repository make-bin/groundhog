// @AI_GENERATED
package vo

import "fmt"

// SecretRef is a value object representing a reference to a secret.
// It is immutable after creation.
type SecretRef struct {
	key    string
	source string
}

// NewSecretRef creates a new SecretRef after validating that key and source are non-empty.
func NewSecretRef(key, source string) (SecretRef, error) {
	if key == "" {
		return SecretRef{}, fmt.Errorf("secret ref key must not be empty")
	}
	if source == "" {
		return SecretRef{}, fmt.Errorf("secret ref source must not be empty")
	}
	return SecretRef{key: key, source: source}, nil
}

// Key returns the secret key.
func (r SecretRef) Key() string { return r.key }

// Source returns the secret source.
func (r SecretRef) Source() string { return r.source }

// @AI_GENERATED: end
