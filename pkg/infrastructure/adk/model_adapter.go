// @AI_GENERATED
package adk

import (
	"context"
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// LLM is the internal interface representing a language model.
type LLM interface {
	// GenerateContent sends a prompt and returns a response.
	GenerateContent(ctx context.Context, prompt string) (string, error)
	// GenerateContentStream sends a prompt and streams chunks to the channel.
	GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error)
}

// ModelFactory is a function that creates an LLM for a given ModelConfig.
type ModelFactory func(ctx context.Context, cfg vo.ModelConfig) (LLM, error)

// ModelAdapter maps domain ModelConfig to LLM instances.
type ModelAdapter struct {
	providers map[vo.ProviderType]ModelFactory
}

// NewModelAdapter creates a ModelAdapter with no registered providers.
func NewModelAdapter() *ModelAdapter {
	return &ModelAdapter{providers: make(map[vo.ProviderType]ModelFactory)}
}

// RegisterProvider registers a factory for the given provider type.
func (a *ModelAdapter) RegisterProvider(provider vo.ProviderType, factory ModelFactory) {
	a.providers[provider] = factory
}

// ToADKModel returns an LLM for the given ModelConfig.
func (a *ModelAdapter) ToADKModel(ctx context.Context, cfg vo.ModelConfig) (LLM, error) {
	factory, ok := a.providers[cfg.Provider()]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider())
	}
	return factory(ctx, cfg)
}

// ToFallbackLLM creates a FallbackLLM for the given ModelConfig.
// It builds the primary LLM and all fallback LLMs from the FallbackChain.
// If no fallbacks are configured, returns the primary LLM directly.
func (a *ModelAdapter) ToFallbackLLM(ctx context.Context, cfg vo.ModelConfig, log logger.Logger) (LLM, error) {
	primary, err := a.ToADKModel(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create primary model: %w", err)
	}

	chain := cfg.FallbackChain()
	if len(chain) == 0 {
		return primary, nil
	}

	fallbacks := make([]LLM, 0, len(chain))
	for _, modelName := range chain {
		fallbackCfg, cfgErr := vo.NewModelConfig(
			cfg.Provider(),
			modelName,
			cfg.Temperature(),
			cfg.MaxTokens(),
			nil,
			cfg.AuthProfile(),
		)
		if cfgErr != nil {
			log.Warn("failed to create fallback model config", "model", modelName, "error", cfgErr)
			continue
		}
		fb, fbErr := a.ToADKModel(ctx, fallbackCfg)
		if fbErr != nil {
			log.Warn("failed to create fallback model", "model", modelName, "error", fbErr)
			continue
		}
		fallbacks = append(fallbacks, fb)
	}

	if len(fallbacks) == 0 {
		return primary, nil
	}

	return NewFallbackLLM(primary, fallbacks, log), nil
}

// LLMToolCall represents a single tool call request returned by the LLM.
type LLMToolCall struct {
	ID        string // tool call id, used to correlate tool_result
	Name      string // function name
	Arguments string // JSON string of arguments
}

// LLMResponse is the structured response from the LLM, containing either text content or tool calls.
type LLMResponse struct {
	Content   string        // text reply (non-empty when ToolCalls is empty)
	ToolCalls []LLMToolCall // tool call requests (non-empty when LLM wants to call tools)
}

// LLMMessage represents a single message in the conversation history.
type LLMMessage struct {
	Role       string        // "system" | "user" | "assistant" | "tool"
	Content    string        // message content
	ToolCallID string        // for role="tool": the ID of the tool call being responded to
	ToolCalls  []LLMToolCall // for role="assistant" with tool calls
}

// ToolSchemaFunc holds the function definition within a ToolSchema.
type ToolSchemaFunc struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// ToolSchema is the tool definition format sent to the LLM (OpenAI function calling format).
type ToolSchema struct {
	Type     string         `json:"type"` // always "function"
	Function ToolSchemaFunc `json:"function"`
}

// LLMWithTools extends the LLM interface to support structured tool use.
type LLMWithTools interface {
	LLM
	// ChatWithTools sends a structured message list with tool definitions and returns a structured response.
	ChatWithTools(ctx context.Context, messages []LLMMessage, tools []ToolSchema) (LLMResponse, error)
	// ChatWithToolsStream is like ChatWithTools but streams text chunks when the response is pure text.
	// If the response contains tool calls, the channel receives a single empty string and closes.
	// The second return value is the full LLMResponse (for tool calls); for text responses it is empty.
	ChatWithToolsStream(ctx context.Context, messages []LLMMessage, tools []ToolSchema) (<-chan string, *LLMResponse, error)
}

// @AI_GENERATED: end
