// @AI_GENERATED
package vo

import "fmt"

// CredentialType represents the type of credential.
type CredentialType int

const (
	CredentialPassword CredentialType = iota
	CredentialToken
	CredentialOAuth
	CredentialAPIKey
)

// Credential is a value object representing an authentication credential.
// It is immutable after creation.
type Credential struct {
	credType CredentialType
	secret   string
	issuer   string
}

// NewCredential creates a new Credential after validating that the secret is non-empty.
func NewCredential(t CredentialType, secret, issuer string) (Credential, error) {
	if secret == "" {
		return Credential{}, fmt.Errorf("credential secret must not be empty")
	}
	return Credential{credType: t, secret: secret, issuer: issuer}, nil
}

// Type returns the credential type.
func (c Credential) Type() CredentialType { return c.credType }

// Secret returns the credential secret.
func (c Credential) Secret() string { return c.secret }

// Issuer returns the credential issuer.
func (c Credential) Issuer() string { return c.issuer }

// @AI_GENERATED: end
