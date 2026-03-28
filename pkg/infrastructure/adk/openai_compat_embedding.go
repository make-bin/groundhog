package adk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	memory_service "github.com/make-bin/groundhog/pkg/domain/memory/service"
)

type openAICompatEmbeddingProvider struct {
	baseURL    string
	apiKey     string
	modelName  string
	httpClient *http.Client
}

// NewOpenAICompatEmbeddingProvider creates a new EmbeddingProvider using OpenAI-compatible API.
func NewOpenAICompatEmbeddingProvider(baseURL, apiKey, modelName string) memory_service.EmbeddingProvider {
	return &openAICompatEmbeddingProvider{
		baseURL:    baseURL,
		apiKey:     apiKey,
		modelName:  modelName,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type embeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func (p *openAICompatEmbeddingProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	reqBody, err := json.Marshal(embeddingRequest{Model: p.modelName, Input: text})
	if err != nil {
		return nil, fmt.Errorf("embedding: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/embeddings", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("embedding: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var embResp embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
		return nil, fmt.Errorf("embedding: decode response: %w", err)
	}
	if len(embResp.Data) == 0 {
		return nil, fmt.Errorf("embedding: empty response data")
	}

	f64 := embResp.Data[0].Embedding
	result := make([]float32, len(f64))
	for i, v := range f64 {
		result[i] = float32(v)
	}
	return result, nil
}
