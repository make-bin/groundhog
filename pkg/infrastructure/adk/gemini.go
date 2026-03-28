// @AI_GENERATED
package adk

import (
	"context"
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

type geminiLLM struct {
	modelName string
}

func (g *geminiLLM) GenerateContent(ctx context.Context, prompt string) (string, error) {
	return fmt.Sprintf("[gemini:%s] %s", g.modelName, prompt), nil
}

func (g *geminiLLM) GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error) {
	ch := make(chan string, 1)
	go func() {
		defer close(ch)
		ch <- fmt.Sprintf("[gemini:%s] %s", g.modelName, prompt)
	}()
	return ch, nil
}

// GeminiProviderFactory returns a ModelFactory for Gemini models.
func GeminiProviderFactory() ModelFactory {
	return func(ctx context.Context, cfg vo.ModelConfig) (LLM, error) {
		return &geminiLLM{modelName: cfg.ModelName()}, nil
	}
}

// @AI_GENERATED: end
