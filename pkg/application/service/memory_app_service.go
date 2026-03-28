package service

import (
	"context"
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/application/dto"
	memorydomain "github.com/make-bin/groundhog/pkg/domain/memory"
	memory "github.com/make-bin/groundhog/pkg/domain/memory/aggregate/memory"
	memory_repository "github.com/make-bin/groundhog/pkg/domain/memory/repository"
	memory_service "github.com/make-bin/groundhog/pkg/domain/memory/service"
	"github.com/make-bin/groundhog/pkg/domain/memory/vo"
)

// MemoryAppService defines the application service interface for memory management.
type MemoryAppService interface {
	SaveMemory(ctx context.Context, userID, content string) (*dto.MemoryResponse, error)
	SearchMemory(ctx context.Context, userID, query string, limit int) ([]*dto.MemorySearchResponse, error)
	ListMemories(ctx context.Context, userID string, offset, limit int) (*dto.MemoryListResponse, error)
	GetMemory(ctx context.Context, id, userID string) (*dto.MemoryResponse, error)
	UpdateMemory(ctx context.Context, id, userID, content string) (*dto.MemoryResponse, error)
	DeleteMemory(ctx context.Context, id, userID string) error
}

type memoryAppService struct {
	MemoryRepo        memory_repository.MemoryRepository `inject:""`
	EmbeddingProvider memory_service.EmbeddingProvider   `inject:"embedding_provider"`
}

// NewMemoryAppService creates a new MemoryAppService. Dependencies are injected via struct tags.
func NewMemoryAppService() MemoryAppService {
	return &memoryAppService{}
}

func (s *memoryAppService) SaveMemory(ctx context.Context, userID, content string) (*dto.MemoryResponse, error) {
	embedding, err := s.EmbeddingProvider.GenerateEmbedding(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("memory_app_service: generate embedding: %w", err)
	}

	memID, err := vo.NewMemoryID(fmt.Sprintf("mem-%d", time.Now().UnixNano()))
	if err != nil {
		return nil, err
	}

	m, err := memory.NewMemory(memID, userID, content, embedding)
	if err != nil {
		return nil, err
	}

	if err := s.MemoryRepo.Create(ctx, m); err != nil {
		return nil, err
	}

	return toMemoryResponse(m), nil
}

func (s *memoryAppService) SearchMemory(ctx context.Context, userID, query string, limit int) ([]*dto.MemorySearchResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	embedding, err := s.EmbeddingProvider.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("memory_app_service: generate embedding: %w", err)
	}

	results, err := s.MemoryRepo.HybridSearch(ctx, userID, embedding, query, limit, 0.7, 0.3)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.MemorySearchResponse, 0, len(results))
	for _, r := range results {
		responses = append(responses, &dto.MemorySearchResponse{
			Memory:      *toMemoryResponse(r.Memory),
			HybridScore: r.HybridScore,
		})
	}
	return responses, nil
}

func (s *memoryAppService) ListMemories(ctx context.Context, userID string, offset, limit int) (*dto.MemoryListResponse, error) {
	filter := memory_repository.MemoryFilter{
		UserID: &userID,
		Offset: offset,
		Limit:  limit,
	}

	memories, total, err := s.MemoryRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.MemoryResponse, 0, len(memories))
	for _, m := range memories {
		responses = append(responses, toMemoryResponse(m))
	}

	return &dto.MemoryListResponse{
		Memories: responses,
		Total:    total,
		Offset:   offset,
		Limit:    limit,
	}, nil
}

func (s *memoryAppService) GetMemory(ctx context.Context, id, userID string) (*dto.MemoryResponse, error) {
	memID, err := vo.NewMemoryID(id)
	if err != nil {
		return nil, err
	}

	m, err := s.MemoryRepo.FindByID(ctx, memID, userID)
	if err != nil {
		return nil, err
	}

	return toMemoryResponse(m), nil
}

func (s *memoryAppService) UpdateMemory(ctx context.Context, id, userID, content string) (*dto.MemoryResponse, error) {
	memID, err := vo.NewMemoryID(id)
	if err != nil {
		return nil, err
	}

	m, err := s.MemoryRepo.FindByID(ctx, memID, userID)
	if err != nil {
		return nil, err
	}

	embedding, err := s.EmbeddingProvider.GenerateEmbedding(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("memory_app_service: generate embedding: %w", err)
	}

	m.UpdateContent(content, embedding)

	if err := s.MemoryRepo.Update(ctx, m); err != nil {
		return nil, err
	}

	return toMemoryResponse(m), nil
}

func (s *memoryAppService) DeleteMemory(ctx context.Context, id, userID string) error {
	memID, err := vo.NewMemoryID(id)
	if err != nil {
		return err
	}

	m, err := s.MemoryRepo.FindByID(ctx, memID, userID)
	if err != nil {
		return err
	}

	if m.UserID() != userID {
		return memorydomain.ErrMemoryAccessDenied
	}

	return s.MemoryRepo.Delete(ctx, memID, userID)
}

func toMemoryResponse(m *memory.Memory) *dto.MemoryResponse {
	return &dto.MemoryResponse{
		ID:        m.ID().Value(),
		UserID:    m.UserID(),
		Content:   m.Content(),
		Tags:      m.Tags(),
		CreatedAt: m.CreatedAt(),
		UpdatedAt: m.UpdatedAt(),
	}
}
