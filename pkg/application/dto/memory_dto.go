package dto

import "time"

// MemoryResponse is the response DTO for a memory.
type MemoryResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMemoryRequest is the request DTO for creating a memory.
type CreateMemoryRequest struct {
	Content string   `json:"content" binding:"required,min=1"`
	Tags    []string `json:"tags"`
}

// UpdateMemoryRequest is the request DTO for updating a memory.
type UpdateMemoryRequest struct {
	Content string   `json:"content" binding:"required,min=1"`
	Tags    []string `json:"tags"`
}

// MemorySearchRequest is the request DTO for searching memories.
type MemorySearchRequest struct {
	Query string `json:"query" binding:"required"`
	Limit int    `json:"limit"`
}

// MemorySearchResponse is the response DTO for a memory search result.
type MemorySearchResponse struct {
	Memory      MemoryResponse `json:"memory"`
	HybridScore float32        `json:"score"`
}

// MemoryListResponse is the response DTO for a paginated list of memories.
type MemoryListResponse struct {
	Memories []*MemoryResponse `json:"memories"`
	Total    int               `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}
