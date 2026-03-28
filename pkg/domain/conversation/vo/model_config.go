// @AI_GENERATED
package vo

import "fmt"

// ModelConfig represents the configuration for an LLM model.
// It is immutable after creation.
type ModelConfig struct {
	provider      ProviderType
	modelName     string
	temperature   float64
	maxTokens     int
	fallbackChain []string
	authProfile   string
}

// NewModelConfig creates a new ModelConfig after validating required fields.
func NewModelConfig(provider ProviderType, modelName string, temperature float64, maxTokens int, fallbackChain []string, authProfile string) (ModelConfig, error) {
	if modelName == "" {
		return ModelConfig{}, fmt.Errorf("model name must not be empty")
	}
	// Copy fallbackChain to ensure immutability.
	chain := make([]string, len(fallbackChain))
	copy(chain, fallbackChain)
	return ModelConfig{
		provider:      provider,
		modelName:     modelName,
		temperature:   temperature,
		maxTokens:     maxTokens,
		fallbackChain: chain,
		authProfile:   authProfile,
	}, nil
}

// Provider returns the provider type.
func (m ModelConfig) Provider() ProviderType { return m.provider }

// ModelName returns the model name.
func (m ModelConfig) ModelName() string { return m.modelName }

// Temperature returns the temperature setting.
func (m ModelConfig) Temperature() float64 { return m.temperature }

// MaxTokens returns the maximum token limit.
func (m ModelConfig) MaxTokens() int { return m.maxTokens }

// FallbackChain returns a copy of the fallback chain.
func (m ModelConfig) FallbackChain() []string {
	chain := make([]string, len(m.fallbackChain))
	copy(chain, m.fallbackChain)
	return chain
}

// AuthProfile returns the authentication profile.
func (m ModelConfig) AuthProfile() string { return m.authProfile }

// @AI_GENERATED: end
