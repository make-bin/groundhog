// @AI_GENERATED
package service

import "github.com/make-bin/groundhog/pkg/domain/identity/vo"

// RateLimitService defines the domain service interface for rate limiting.
type RateLimitService interface {
	// Allow returns true if the principal is allowed to make a request.
	Allow(principalID vo.PrincipalID) bool
}

// @AI_GENERATED: end
