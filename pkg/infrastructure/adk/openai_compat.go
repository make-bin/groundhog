// @AI_GENERATED
package adk

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

type openAICompatLLM struct {
	baseURL    string
	modelName  string
	apiKey     string
	httpClient *http.Client
}

func (o *openAICompatLLM) GenerateContent(ctx context.Context, prompt string) (string, error) {
	payload := map[string]any{
		"model": o.modelName,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("openai_compat: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("openai_compat: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai_compat: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("openai_compat: read response: %w", err)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("openai_compat: unmarshal response: %w", err)
	}
	if result.Error != nil {
		return "", fmt.Errorf("openai_compat: api error: %s", result.Error.Message)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai_compat: no choices in response")
	}
	return result.Choices[0].Message.Content, nil
}

func (o *openAICompatLLM) GenerateContentStream(ctx context.Context, prompt string) (<-chan string, error) {
	payload := map[string]any{
		"model": o.modelName,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"stream": true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("openai_compat: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai_compat: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai_compat: request failed: %w", err)
	}

	ch := make(chan string, 16)
	go func() {
		defer close(ch)
		defer resp.Body.Close()
		streamOpenAICompatSSE(ctx, resp.Body, ch)
	}()

	return ch, nil
}

// OpenAICompatProviderFactory returns a ModelFactory for OpenAI-compatible APIs.
// baseURL is the API base (e.g. https://api.siliconflow.cn/v1), apiKey is the bearer token.
func OpenAICompatProviderFactory(baseURL, apiKey string) ModelFactory {
	return func(ctx context.Context, cfg vo.ModelConfig) (LLM, error) {
		key := apiKey
		if cfg.AuthProfile() != "" {
			key = cfg.AuthProfile()
		}
		return &openAICompatLLM{
			baseURL:    baseURL,
			modelName:  cfg.ModelName(),
			apiKey:     key,
			httpClient: &http.Client{Timeout: 120 * time.Second},
		}, nil
	}
}

// streamOpenAICompatSSE reads an OpenAI-compatible SSE stream and sends
// content delta tokens to ch. It stops on [DONE] or context cancellation.
func streamOpenAICompatSSE(ctx context.Context, body io.Reader, ch chan<- string) {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if strings.TrimSpace(data) == "[DONE]" {
			return
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			select {
			case ch <- chunk.Choices[0].Delta.Content:
			case <-ctx.Done():
				return
			}
		}
	}
}

// @AI_GENERATED: end

// ChatWithTools implements LLMWithTools. It sends a structured message list with tool definitions
// and returns a structured response containing either text content or tool call requests.
func (o *openAICompatLLM) ChatWithTools(ctx context.Context, messages []LLMMessage, tools []ToolSchema) (LLMResponse, error) {
	// Build the messages array for the request
	type reqToolCall struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Function struct {
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
		} `json:"function"`
	}
	type reqMessage struct {
		Role       string        `json:"role"`
		Content    string        `json:"content,omitempty"`
		ToolCallID string        `json:"tool_call_id,omitempty"`
		ToolCalls  []reqToolCall `json:"tool_calls,omitempty"`
	}

	reqMsgs := make([]reqMessage, 0, len(messages))
	for _, m := range messages {
		rm := reqMessage{Role: m.Role, Content: m.Content, ToolCallID: m.ToolCallID}
		if len(m.ToolCalls) > 0 {
			for _, tc := range m.ToolCalls {
				rtc := reqToolCall{ID: tc.ID, Type: "function"}
				rtc.Function.Name = tc.Name
				rtc.Function.Arguments = tc.Arguments
				rm.ToolCalls = append(rm.ToolCalls, rtc)
			}
		}
		reqMsgs = append(reqMsgs, rm)
	}

	payload := map[string]any{
		"model":    o.modelName,
		"messages": reqMsgs,
	}
	if len(tools) > 0 {
		payload["tools"] = tools
		payload["tool_choice"] = "auto"
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("openai_compat: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return LLMResponse{}, fmt.Errorf("openai_compat: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("openai_compat: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("openai_compat: read response: %w", err)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
				ToolCalls []struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return LLMResponse{}, fmt.Errorf("openai_compat: unmarshal response: %w", err)
	}
	if result.Error != nil {
		return LLMResponse{}, fmt.Errorf("openai_compat: api error: %s", result.Error.Message)
	}
	if len(result.Choices) == 0 {
		return LLMResponse{}, fmt.Errorf("openai_compat: no choices in response")
	}

	msg := result.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		toolCalls := make([]LLMToolCall, 0, len(msg.ToolCalls))
		for _, tc := range msg.ToolCalls {
			toolCalls = append(toolCalls, LLMToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			})
		}
		return LLMResponse{ToolCalls: toolCalls}, nil
	}
	return LLMResponse{Content: msg.Content}, nil
}

// ChatWithToolsStream streams the LLM response when it is pure text, or returns tool calls non-streaming.
// Strategy: send stream=true; if the first meaningful delta contains tool_calls, fall back to non-streaming
// by collecting the full response and returning it via LLMResponse. Otherwise stream text chunks.
func (o *openAICompatLLM) ChatWithToolsStream(ctx context.Context, messages []LLMMessage, tools []ToolSchema) (<-chan string, *LLMResponse, error) {
	type reqToolCall struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Function struct {
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
		} `json:"function"`
	}
	type reqMessage struct {
		Role       string        `json:"role"`
		Content    string        `json:"content,omitempty"`
		ToolCallID string        `json:"tool_call_id,omitempty"`
		ToolCalls  []reqToolCall `json:"tool_calls,omitempty"`
	}

	reqMsgs := make([]reqMessage, 0, len(messages))
	for _, m := range messages {
		rm := reqMessage{Role: m.Role, Content: m.Content, ToolCallID: m.ToolCallID}
		for _, tc := range m.ToolCalls {
			rtc := reqToolCall{ID: tc.ID, Type: "function"}
			rtc.Function.Name = tc.Name
			rtc.Function.Arguments = tc.Arguments
			rm.ToolCalls = append(rm.ToolCalls, rtc)
		}
		reqMsgs = append(reqMsgs, rm)
	}

	payload := map[string]any{
		"model":    o.modelName,
		"messages": reqMsgs,
		"stream":   true,
	}
	if len(tools) > 0 {
		payload["tools"] = tools
		payload["tool_choice"] = "auto"
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("openai_compat stream: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, nil, fmt.Errorf("openai_compat stream: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("openai_compat stream: request: %w", err)
	}

	ch := make(chan string, 32)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		// Streaming delta structs
		type deltaToolCallFunc struct {
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
		}
		type deltaToolCall struct {
			Index    int               `json:"index"`
			ID       string            `json:"id"`
			Type     string            `json:"type"`
			Function deltaToolCallFunc `json:"function"`
		}
		type delta struct {
			Content   string          `json:"content"`
			ToolCalls []deltaToolCall `json:"tool_calls"`
		}
		type choice struct {
			Delta        delta  `json:"delta"`
			FinishReason string `json:"finish_reason"`
		}
		type streamChunk struct {
			Choices []choice `json:"choices"`
		}

		// Accumulate tool call fragments indexed by tool call index
		type tcAccum struct {
			id        string
			name      string
			arguments strings.Builder
		}
		tcMap := map[int]*tcAccum{}
		hasToolCalls := false

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if strings.TrimSpace(data) == "[DONE]" {
				break
			}

			var sc streamChunk
			if err := json.Unmarshal([]byte(data), &sc); err != nil || len(sc.Choices) == 0 {
				continue
			}

			d := sc.Choices[0].Delta

			// Accumulate tool call deltas
			if len(d.ToolCalls) > 0 {
				hasToolCalls = true
				for _, tc := range d.ToolCalls {
					acc, ok := tcMap[tc.Index]
					if !ok {
						acc = &tcAccum{}
						tcMap[tc.Index] = acc
					}
					if tc.ID != "" {
						acc.id = tc.ID
					}
					if tc.Function.Name != "" {
						acc.name += tc.Function.Name
					}
					acc.arguments.WriteString(tc.Function.Arguments)
				}
				continue
			}

			// Stream text content
			if d.Content != "" && !hasToolCalls {
				select {
				case ch <- d.Content:
				case <-ctx.Done():
					return
				}
			}
		}

		// If tool calls were detected, send them as a special marker via a side channel.
		// We encode the LLMResponse as JSON and send it as a single sentinel chunk prefixed with "\x00tc:".
		if hasToolCalls {
			toolCalls := make([]LLMToolCall, 0, len(tcMap))
			for i := 0; i < len(tcMap); i++ {
				acc, ok := tcMap[i]
				if !ok {
					continue
				}
				toolCalls = append(toolCalls, LLMToolCall{
					ID:        acc.id,
					Name:      acc.name,
					Arguments: acc.arguments.String(),
				})
			}
			encoded, _ := json.Marshal(toolCalls)
			select {
			case ch <- "\x00tc:" + string(encoded):
			case <-ctx.Done():
			}
		}
	}()

	return ch, nil, nil
}
