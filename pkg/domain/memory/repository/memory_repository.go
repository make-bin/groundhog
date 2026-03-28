package repository

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/memory/aggregate/memory"
	"github.com/make-bin/groundhog/pkg/domain/memory/vo"
)

// MemoryFilter defines the filter criteria for listing memories.
type MemoryFilter struct {
	UserID *string
	Limit  int
	Offset int
}

// MemorySearchResult holds a memory along with its search scores.
type MemorySearchResult struct {
	Memory       *memory.Memory
	VectorScore  float32
	KeywordScore float32
	HybridScore  float32
}

// MemoryRepository defines the data access contract for the Memory aggregate.
type MemoryRepository interface {
	Create(ctx context.Context, m *memory.Memory) error
	FindByID(ctx context.Context, id vo.MemoryID, userID string) (*memory.Memory, error)
	Update(ctx context.Context, m *memory.Memory) error
	Delete(ctx context.Context, id vo.MemoryID, userID string) error
	List(ctx context.Context, filter MemoryFilter) ([]*memory.Memory, int, error)
	HybridSearch(ctx context.Context, userID string, embedding []float32, keyword string, limit int, vectorWeight, keywordWeight float32) ([]*MemorySearchResult, error)
}
