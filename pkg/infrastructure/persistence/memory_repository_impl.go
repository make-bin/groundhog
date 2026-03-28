package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	memorydomain "github.com/make-bin/groundhog/pkg/domain/memory"
	"github.com/make-bin/groundhog/pkg/domain/memory/aggregate/memory"
	"github.com/make-bin/groundhog/pkg/domain/memory/repository"
	"github.com/make-bin/groundhog/pkg/domain/memory/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/mapper"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

type memoryRepositoryImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
}

// NewMemoryRepository creates a new MemoryRepository implementation.
func NewMemoryRepository() repository.MemoryRepository {
	return &memoryRepositoryImpl{}
}

// Create persists a new Memory aggregate.
func (r *memoryRepositoryImpl) Create(ctx context.Context, m *memory.Memory) error {
	memPO := mapper.DomainToPO(m)
	return r.DataStore.DB().WithContext(ctx).Create(memPO).Error
}

// FindByID retrieves a Memory aggregate by its MemoryID and userID.
func (r *memoryRepositoryImpl) FindByID(ctx context.Context, id vo.MemoryID, userID string) (*memory.Memory, error) {
	var memPO po.MemoryPO
	result := r.DataStore.DB().WithContext(ctx).
		Where("memory_id = ? AND user_id = ? AND deleted_at IS NULL", id.Value(), userID).
		First(&memPO)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, memorydomain.ErrMemoryNotFound
		}
		return nil, result.Error
	}
	return mapper.POToDomain(&memPO)
}

// Update persists changes to an existing Memory aggregate.
func (r *memoryRepositoryImpl) Update(ctx context.Context, m *memory.Memory) error {
	return r.DataStore.DB().WithContext(ctx).
		Model(&po.MemoryPO{}).
		Where("memory_id = ? AND deleted_at IS NULL", m.ID().Value()).
		Updates(map[string]interface{}{
			"content":    m.Content(),
			"embedding":  mapper.EmbeddingToText(m.Embedding()),
			"updated_at": m.UpdatedAt(),
		}).Error
}

// Delete soft-deletes a Memory aggregate by its MemoryID and userID.
func (r *memoryRepositoryImpl) Delete(ctx context.Context, id vo.MemoryID, userID string) error {
	return r.DataStore.DB().WithContext(ctx).
		Where("memory_id = ? AND user_id = ?", id.Value(), userID).
		Delete(&po.MemoryPO{}).Error
}

// List retrieves Memory aggregates matching the filter with pagination.
func (r *memoryRepositoryImpl) List(ctx context.Context, filter repository.MemoryFilter) ([]*memory.Memory, int, error) {
	db := r.DataStore.DB().WithContext(ctx).Model(&po.MemoryPO{}).Where("deleted_at IS NULL")
	if filter.UserID != nil {
		db = db.Where("user_id = ?", *filter.UserID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	var memPOs []*po.MemoryPO
	if err := db.Offset(filter.Offset).Limit(limit).Find(&memPOs).Error; err != nil {
		return nil, 0, err
	}

	memories, err := mapper.POListToDomainList(memPOs)
	if err != nil {
		return nil, 0, err
	}
	return memories, int(total), nil
}

// hybridSearchRow is a local struct for scanning HybridSearch raw SQL results.
type hybridSearchRow struct {
	po.MemoryPO
	HybridScore  float32 `gorm:"column:hybrid_score"`
	VectorScore  float32 `gorm:"column:vector_score"`
	KeywordScore float32 `gorm:"column:keyword_score"`
}

// HybridSearch performs a combined vector + full-text search.
// Falls back to full-text only when pgvector is not available.
func (r *memoryRepositoryImpl) HybridSearch(ctx context.Context, userID string, embedding []float32, keyword string, limit int, vectorWeight, keywordWeight float32) ([]*repository.MemorySearchResult, error) {
	if r.pgvectorAvailable(ctx) {
		return r.hybridSearchVector(ctx, userID, embedding, keyword, limit, vectorWeight, keywordWeight)
	}
	return r.fullTextSearch(ctx, userID, keyword, limit, keywordWeight)
}

// pgvectorAvailable checks whether the vector extension is installed.
func (r *memoryRepositoryImpl) pgvectorAvailable(ctx context.Context) bool {
	var count int64
	r.DataStore.DB().WithContext(ctx).
		Raw("SELECT COUNT(*) FROM pg_extension WHERE extname = 'vector'").
		Scan(&count)
	return count > 0
}

// hybridSearchVector uses pgvector cosine distance + ts_rank.
func (r *memoryRepositoryImpl) hybridSearchVector(ctx context.Context, userID string, embedding []float32, keyword string, limit int, vectorWeight, keywordWeight float32) ([]*repository.MemorySearchResult, error) {
	vecStr := float32SliceToVectorString(embedding)

	rawSQL := `
		SELECT *,
			(1 - (embedding <=> ?::vector)) * ? AS vector_score,
			ts_rank(to_tsvector('simple', content), plainto_tsquery('simple', ?)) * ? AS keyword_score,
			(1 - (embedding <=> ?::vector)) * ? + ts_rank(to_tsvector('simple', content), plainto_tsquery('simple', ?)) * ? AS hybrid_score
		FROM memories
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY hybrid_score DESC
		LIMIT ?`

	var rows []hybridSearchRow
	if err := r.DataStore.DB().WithContext(ctx).Raw(rawSQL,
		vecStr, vectorWeight,
		keyword, keywordWeight,
		vecStr, vectorWeight, keyword, keywordWeight,
		userID, limit,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return r.rowsToResults(rows)
}

// fullTextSearch falls back to ts_rank only when pgvector is unavailable.
func (r *memoryRepositoryImpl) fullTextSearch(ctx context.Context, userID string, keyword string, limit int, keywordWeight float32) ([]*repository.MemorySearchResult, error) {
	rawSQL := `
		SELECT *,
			0.0 AS vector_score,
			ts_rank(to_tsvector('simple', content), plainto_tsquery('simple', ?)) * ? AS keyword_score,
			ts_rank(to_tsvector('simple', content), plainto_tsquery('simple', ?)) * ? AS hybrid_score
		FROM memories
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY hybrid_score DESC
		LIMIT ?`

	var rows []hybridSearchRow
	if err := r.DataStore.DB().WithContext(ctx).Raw(rawSQL,
		keyword, keywordWeight,
		keyword, keywordWeight,
		userID, limit,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return r.rowsToResults(rows)
}

func (r *memoryRepositoryImpl) rowsToResults(rows []hybridSearchRow) ([]*repository.MemorySearchResult, error) {
	results := make([]*repository.MemorySearchResult, 0, len(rows))
	for i := range rows {
		m, err := mapper.POToDomain(&rows[i].MemoryPO)
		if err != nil {
			return nil, err
		}
		results = append(results, &repository.MemorySearchResult{
			Memory:       m,
			VectorScore:  rows[i].VectorScore,
			KeywordScore: rows[i].KeywordScore,
			HybridScore:  rows[i].HybridScore,
		})
	}
	return results, nil
}

// float32SliceToVectorString converts a []float32 to a pgvector-compatible string like "[0.1,0.2,0.3]".
func float32SliceToVectorString(v []float32) string {
	parts := make([]string, len(v))
	for i, f := range v {
		parts[i] = fmt.Sprintf("%g", f)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
