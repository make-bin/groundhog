// @AI_GENERATED
package adk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

const mistralAPIBase = "https://api.mistral.ai/v1"

type mistralLLM struct {
	modelName  string
	apiKey     string
	httpClient *http.Client
}

func (m *mistralLLM) GenerateContent(ctx context.Context, prompt string) (string, error) {
	payload := map[string]any{
		"model": m.modelName,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("mistral: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		mistralAPIBase+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("mistral: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("mistral: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("mistral: read response: %w", err)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("mistral: unmarshal response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("mistral: no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}

func (m *mistralLLM) GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error) {
	payload := map[string]any{
		"model": m.modelName,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"stream": true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("mistral: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		mistralAPIBase+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("mistral: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mistral: request failed: %w", err)
	}

	ch := make(chan string, 16)
	go func() {
		defer close(ch)
		defer resp.Body.Close()
		streamOpenAICompatSSE(ctx, resp.Body, ch)
	}()

	return ch, nil
}

// MistralProviderFactory returns a ModelFactory for Mistral models.
func MistralProviderFactory() ModelFactory {
	return func(ctx context.Context, cfg vo.ModelConfig) (LLM, error) {
		return &mistralLLM{
			modelName:  cfg.ModelName(),
			apiKey:     cfg.AuthProfile(),
			httpClient: &http.Client{Timeout: 60 * time.Second},
		}, nil
	}
}

// @AI_GENERATED: end
