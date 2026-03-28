// @AI_GENERATED
package adk

import (
	"errors"
	"sync"
	"time"
)

var ErrNoAvailableKey = errors.New("no available API key for provider")

// APIKeyProfile holds a single API key with usage tracking.
type APIKeyProfile struct {
	Key           string
	UsageCount    int64
	CooldownUntil time.Time
}

// IsAvailable returns true if the key is not in cooldown.
func (p *APIKeyProfile) IsAvailable() bool {
	return time.Now().After(p.CooldownUntil)
}

// AuthProfileManager manages multiple API keys per provider with rotation and cooldown.
type AuthProfileManager struct {
	mu       sync.RWMutex
	profiles map[string][]*APIKeyProfile // provider name → keys
	indices  map[string]int              // provider name → current round-robin index
}

// NewAuthProfileManager creates a new AuthProfileManager.
func NewAuthProfileManager() *AuthProfileManager {
	return &AuthProfileManager{
		profiles: make(map[string][]*APIKeyProfile),
		indices:  make(map[string]int),
	}
}

// AddKey registers an API key for the given provider.
func (m *AuthProfileManager) AddKey(provider, key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.profiles[provider] = append(m.profiles[provider], &APIKeyProfile{Key: key})
}

// GetKey returns the next available API key for the provider using round-robin.
// Returns ErrNoAvailableKey if all keys are in cooldown.
func (m *AuthProfileManager) GetKey(provider string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	keys := m.profiles[provider]
	if len(keys) == 0 {
		return "", ErrNoAvailableKey
	}

	start := m.indices[provider]
	for i := 0; i < len(keys); i++ {
		idx := (start + i) % len(keys)
		k := keys[idx]
		if k.IsAvailable() {
			k.UsageCount++
			m.indices[provider] = (idx + 1) % len(keys)
			return k.Key, nil
		}
	}
	return "", ErrNoAvailableKey
}

// MarkRateLimited marks a key as rate-limited with the given cooldown duration.
func (m *AuthProfileManager) MarkRateLimited(provider, key string, cooldown time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, k := range m.profiles[provider] {
		if k.Key == key {
			k.CooldownUntil = time.Now().Add(cooldown)
			return
		}
	}
}

// UsageCount returns the usage count for a specific key.
func (m *AuthProfileManager) UsageCount(provider, key string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, k := range m.profiles[provider] {
		if k.Key == key {
			return k.UsageCount
		}
	}
	return 0
}

// @AI_GENERATED: end
