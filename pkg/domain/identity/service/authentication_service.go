// @AI_GENERATED
package service

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/identity/aggregate/principal"
	"github.com/make-bin/groundhog/pkg/domain/identity/vo"
)

// AuthenticationService defines the domain service interface for authentication.
type AuthenticationService interface {
	// Authenticate verifies the given credential and returns the associated principal.
	Authenticate(ctx context.Context, cred vo.Credential) (*principal.Principal, error)

	// ValidateToken validates a token string and returns the associated principal.
	ValidateToken(ctx context.Context, token string) (*principal.Principal, error)
}

// @AI_GENERATED: end
