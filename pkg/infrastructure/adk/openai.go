// @AI_GENERATED
package adk

import (
	"context"
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

type openAILLM struct {
	modelName string
}

func (o *openAILLM) GenerateContent(ctx context.Context, prompt string) (string, error) {
	return fmt.Sprintf("[openai:%s] %s", o.modelName, prompt), nil
}

func (o *openAILLM) GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error) {
	ch := make(chan string, 1)
	go func() {
		defer close(ch)
		ch <- fmt.Sprintf("[openai:%s] %s", o.modelName, prompt)
	}()
	return ch, nil
}

// OpenAIProviderFactory returns a ModelFactory for OpenAI models.
func OpenAIProviderFactory() ModelFactory {
	return func(ctx context.Context, cfg vo.ModelConfig) (LLM, error) {
		return &openAILLM{modelName: cfg.ModelName()}, nil
	}
}

// @AI_GENERATED: end
