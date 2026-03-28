package mapper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/make-bin/groundhog/pkg/domain/memory/aggregate/memory"
	"github.com/make-bin/groundhog/pkg/domain/memory/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

// EmbeddingToText serializes []float32 to a pgvector-compatible string "[0.1,0.2,...]".
func EmbeddingToText(v []float32) string {
	return embeddingToText(v)
}

// embeddingToText serializes []float32 to a pgvector-compatible string "[0.1,0.2,...]".
func embeddingToText(v []float32) string {
	if len(v) == 0 {
		return ""
	}
	parts := make([]string, len(v))
	for i, f := range v {
		parts[i] = fmt.Sprintf("%g", f)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// textToEmbedding parses a pgvector literal or JSON array back to []float32.
func textToEmbedding(s string) []float32 {
	if s == "" {
		return nil
	}
	// Strip surrounding brackets if present
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '[' {
		s = s[1 : len(s)-1]
	}
	parts := strings.Split(s, ",")
	result := make([]float32, 0, len(parts))
	for _, p := range parts {
		var f float32
		fmt.Sscanf(strings.TrimSpace(p), "%g", &f)
		result = append(result, f)
	}
	return result
}

// DomainToPO converts a Memory aggregate to a MemoryPO.
func DomainToPO(m *memory.Memory) *po.MemoryPO {
	tagsJSON, _ := json.Marshal(m.Tags())
	return &po.MemoryPO{
		MemoryID:  m.ID().Value(),
		UserID:    m.UserID(),
		Content:   m.Content(),
		Embedding: embeddingToText(m.Embedding()),
		Tags:      string(tagsJSON),
		CreatedAt: m.CreatedAt(),
		UpdatedAt: m.UpdatedAt(),
	}
}

// POToDomain converts a MemoryPO to a Memory aggregate.
func POToDomain(p *po.MemoryPO) (*memory.Memory, error) {
	id, err := vo.NewMemoryID(p.MemoryID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct memory_id: %w", err)
	}
	var tags []string
	if p.Tags != "" {
		if err := json.Unmarshal([]byte(p.Tags), &tags); err != nil {
			return nil, fmt.Errorf("unmarshal tags: %w", err)
		}
	}
	if tags == nil {
		tags = []string{}
	}
	return memory.ReconstructMemory(
		id,
		p.UserID,
		p.Content,
		textToEmbedding(p.Embedding),
		tags,
		p.CreatedAt,
		p.UpdatedAt,
	), nil
}

// POListToDomainList converts a slice of MemoryPO to a slice of Memory aggregates.
func POListToDomainList(pos []*po.MemoryPO) ([]*memory.Memory, error) {
	result := make([]*memory.Memory, 0, len(pos))
	for _, p := range pos {
		m, err := POToDomain(p)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}
