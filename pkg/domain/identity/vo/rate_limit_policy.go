// @AI_GENERATED
package vo

import "fmt"

// RateLimitPolicy is a value object representing a rate limiting policy.
// It is immutable after creation.
type RateLimitPolicy struct {
	requestsPerMinute int
	burstSize         int
}

// NewRateLimitPolicy creates a new RateLimitPolicy after validating that both values are positive.
func NewRateLimitPolicy(rpm, burst int) (RateLimitPolicy, error) {
	if rpm <= 0 {
		return RateLimitPolicy{}, fmt.Errorf("requests per minute must be positive, got %d", rpm)
	}
	if burst <= 0 {
		return RateLimitPolicy{}, fmt.Errorf("burst size must be positive, got %d", burst)
	}
	return RateLimitPolicy{requestsPerMinute: rpm, burstSize: burst}, nil
}

// RequestsPerMinute returns the maximum number of requests allowed per minute.
func (p RateLimitPolicy) RequestsPerMinute() int { return p.requestsPerMinute }

// BurstSize returns the maximum burst size.
func (p RateLimitPolicy) BurstSize() int { return p.burstSize }

// @AI_GENERATED: end
