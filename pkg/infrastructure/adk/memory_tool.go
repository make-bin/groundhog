package adk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/make-bin/groundhog/pkg/application/dto"
	memorydomain "github.com/make-bin/groundhog/pkg/domain/memory"
)

// MemoryService is a local interface satisfied by service.MemoryAppService.
// Defined here to avoid a circular import between infrastructure/adk and application/service.
type MemoryService interface {
	SaveMemory(ctx context.Context, userID, content string) (*dto.MemoryResponse, error)
	SearchMemory(ctx context.Context, userID, query string, limit int) ([]*dto.MemorySearchResponse, error)
}

type memorySaveTool struct {
	svc MemoryService
}

func (t *memorySaveTool) Name() string { return "memory_save" }

func (t *memorySaveTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	userID, ok := MemoryUserIDFromContext(ctx)
	if !ok {
		return "", memorydomain.ErrMissingUserID
	}
	content, _ := args["content"].(string)
	if content == "" {
		return "", fmt.Errorf("memory_save: 'content' argument is required")
	}
	resp, err := t.svc.SaveMemory(ctx, userID, content)
	if err != nil {
		return "", fmt.Errorf("memory_save: %w", err)
	}
	out, _ := json.Marshal(resp)
	return string(out), nil
}

type memorySearchTool struct {
	svc MemoryService
}

func (t *memorySearchTool) Name() string { return "memory_search" }

func (t *memorySearchTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	userID, ok := MemoryUserIDFromContext(ctx)
	if !ok {
		return "", memorydomain.ErrMissingUserID
	}
	query, _ := args["query"].(string)
	if query == "" {
		return "", fmt.Errorf("memory_search: 'query' argument is required")
	}
	limit := 10
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}
	results, err := t.svc.SearchMemory(ctx, userID, query, limit)
	if err != nil {
		// Return empty results on search failure to avoid breaking conversation
		return "[]", nil
	}
	out, _ := json.Marshal(results)
	return string(out), nil
}

// NewMemoryTools creates the memory_save and memory_search tools.
// svc must implement MemoryService (satisfied by service.MemoryAppService).
func NewMemoryTools(svc MemoryService) (Tool, Tool) {
	return &memorySaveTool{svc: svc}, &memorySearchTool{svc: svc}
}

// MemoryToolSchemas returns the ToolSchema definitions for memory_save and memory_search.
func MemoryToolSchemas() []ToolSchema {
	return []ToolSchema{
		{
			Type: "function",
			Function: ToolSchemaFunc{
				Name:        "memory_save",
				Description: "Save important information to long-term memory",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"content": map[string]any{
							"type":        "string",
							"description": "The information to save",
						},
					},
					"required": []string{"content"},
				},
			},
		},
		{
			Type: "function",
			Function: ToolSchemaFunc{
				Name:        "memory_search",
				Description: "Search long-term memory for relevant information",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query",
						},
						"limit": map[string]any{
							"type":        "integer",
							"description": "Max results to return",
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}
}
