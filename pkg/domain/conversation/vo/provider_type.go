// @AI_GENERATED
package vo

// ProviderType represents the type of LLM provider.
type ProviderType int

const (
	ProviderGemini ProviderType = iota
	ProviderOpenAI
	ProviderAnthropic
	ProviderOllama
	ProviderGroq
	ProviderOpenAICompat
)

// String returns the string representation of the provider type.
func (p ProviderType) String() string {
	switch p {
	case ProviderGemini:
		return "gemini"
	case ProviderOpenAI:
		return "openai"
	case ProviderAnthropic:
		return "anthropic"
	case ProviderOllama:
		return "ollama"
	case ProviderGroq:
		return "groq"
	case ProviderOpenAICompat:
		return "openai_compat"
	default:
		return "unknown"
	}
}

// @AI_GENERATED: end
