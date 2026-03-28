// @AI_GENERATED
package adk

import (
	"context"
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

type anthropicLLM struct {
	modelName string
}

func (a *anthropicLLM) GenerateContent(ctx context.Context, prompt string) (string, error) {
	return fmt.Sprintf("[anthropic:%s] %s", a.modelName, prompt), nil
}

func (a *anthropicLLM) GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error) {
	ch := make(chan string, 1)
	go func() {
		defer close(ch)
		ch <- fmt.Sprintf("[anthropic:%s] %s", a.modelName, prompt)
	}()
	return ch, nil
}

// AnthropicProviderFactory returns a ModelFactory for Anthropic models.
func AnthropicProviderFactory() ModelFactory {
	return func(ctx context.Context, cfg vo.ModelConfig) (LLM, error) {
		return &anthropicLLM{modelName: cfg.ModelName()}, nil
	}
}

// @AI_GENERATED: end
