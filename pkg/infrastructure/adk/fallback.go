// @AI_GENERATED
package adk

import (
	"context"
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// ModelFallbackTriggeredEvent is published when a model fallback occurs.
type ModelFallbackTriggeredEvent struct {
	FailedModel   string
	FallbackModel string
	Reason        string
	OccurredAt    time.Time
}

// FallbackLLM wraps a primary LLM with a chain of fallback LLMs.
// When the primary fails, it tries each fallback in order.
type FallbackLLM struct {
	primary    LLM
	fallbacks  []LLM
	logger     logger.Logger
	onFallback func(event ModelFallbackTriggeredEvent)
}

// NewFallbackLLM creates a FallbackLLM with the given primary and fallback LLMs.
func NewFallbackLLM(primary LLM, fallbacks []LLM, log logger.Logger) *FallbackLLM {
	return &FallbackLLM{
		primary:   primary,
		fallbacks: fallbacks,
		logger:    log,
	}
}

// SetFallbackHandler sets a callback for when fallback is triggered.
func (f *FallbackLLM) SetFallbackHandler(fn func(event ModelFallbackTriggeredEvent)) {
	f.onFallback = fn
}

// GenerateContent tries the primary LLM, then each fallback in order.
func (f *FallbackLLM) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := f.primary.GenerateContent(ctx, prompt)
	if err == nil {
		return resp, nil
	}

	f.logger.Warn("primary model failed, trying fallbacks", "error", err)

	for i, fb := range f.fallbacks {
		f.logger.Info("trying fallback model", "index", i)

		if f.onFallback != nil {
			f.onFallback(ModelFallbackTriggeredEvent{
				FailedModel:   fmt.Sprintf("fallback-%d", i-1),
				FallbackModel: fmt.Sprintf("fallback-%d", i),
				Reason:        err.Error(),
				OccurredAt:    time.Now(),
			})
		}

		resp, fbErr := fb.GenerateContent(ctx, prompt)
		if fbErr == nil {
			return resp, nil
		}
		f.logger.Warn("fallback model failed", "index", i, "error", fbErr)
		err = fbErr
	}

	return "", fmt.Errorf("all models failed, last error: %w", err)
}

// GenerateContentStream tries the primary LLM stream, then each fallback.
func (f *FallbackLLM) GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error) {
	ch, err := f.primary.GenerateContentStream(ctx, prompt)
	if err == nil {
		return ch, nil
	}

	f.logger.Warn("primary model stream failed, trying fallbacks", "error", err)

	for i, fb := range f.fallbacks {
		if f.onFallback != nil {
			f.onFallback(ModelFallbackTriggeredEvent{
				FailedModel:   fmt.Sprintf("fallback-%d", i-1),
				FallbackModel: fmt.Sprintf("fallback-%d", i),
				Reason:        err.Error(),
				OccurredAt:    time.Now(),
			})
		}

		ch, fbErr := fb.GenerateContentStream(ctx, prompt)
		if fbErr == nil {
			return ch, nil
		}
		f.logger.Warn("fallback model stream failed", "index", i, "error", fbErr)
		err = fbErr
	}

	return nil, fmt.Errorf("all models failed for streaming, last error: %w", err)
}

// @AI_GENERATED: end
