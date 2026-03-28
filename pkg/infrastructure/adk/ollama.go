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

const ollamaAPIBase = "http://localhost:11434"

type ollamaLLM struct {
	modelName  string
	baseURL    string
	httpClient *http.Client
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (o *ollamaLLM) GenerateContent(ctx context.Context, prompt string) (string, error) {
	payload := ollamaGenerateRequest{
		Model:  o.modelName,
		Prompt: prompt,
		Stream: false,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("ollama: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		o.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ollama: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ollama: read response: %w", err)
	}

	var result ollamaGenerateResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("ollama: unmarshal response: %w", err)
	}

	return result.Response, nil
}

func (o *ollamaLLM) GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error) {
	payload := ollamaGenerateRequest{
		Model:  o.modelName,
		Prompt: prompt,
		Stream: true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ollama: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		o.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ollama: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama: request failed: %w", err)
	}

	ch := make(chan string, 16)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk ollamaGenerateResponse
			if err := decoder.Decode(&chunk); err != nil {
				return
			}
			if chunk.Response != "" {
				select {
				case ch <- chunk.Response:
				case <-ctx.Done():
					return
				}
			}
			if chunk.Done {
				return
			}
		}
	}()

	return ch, nil
}

// OllamaProviderFactory returns a ModelFactory for Ollama local models.
func OllamaProviderFactory() ModelFactory {
	return func(ctx context.Context, cfg vo.ModelConfig) (LLM, error) {
		baseURL := ollamaAPIBase
		if cfg.AuthProfile() != "" {
			baseURL = cfg.AuthProfile()
		}
		return &ollamaLLM{
			modelName:  cfg.ModelName(),
			baseURL:    baseURL,
			httpClient: &http.Client{Timeout: 120 * time.Second},
		}, nil
	}
}

// @AI_GENERATED: end
