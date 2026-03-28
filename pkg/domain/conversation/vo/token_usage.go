// @AI_GENERATED
package vo

// TokenUsage represents the token consumption for a model invocation.
// It is immutable after creation. totalTokens is auto-calculated.
type TokenUsage struct {
	promptTokens     int
	completionTokens int
	totalTokens      int
}

// NewTokenUsage creates a new TokenUsage. totalTokens is automatically
// calculated as promptTokens + completionTokens.
func NewTokenUsage(promptTokens, completionTokens int) TokenUsage {
	return TokenUsage{
		promptTokens:     promptTokens,
		completionTokens: completionTokens,
		totalTokens:      promptTokens + completionTokens,
	}
}

// PromptTokens returns the number of prompt tokens.
func (t TokenUsage) PromptTokens() int { return t.promptTokens }

// CompletionTokens returns the number of completion tokens.
func (t TokenUsage) CompletionTokens() int { return t.completionTokens }

// TotalTokens returns the total number of tokens.
func (t TokenUsage) TotalTokens() int { return t.totalTokens }

// @AI_GENERATED: end
