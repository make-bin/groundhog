package memory

import (
	"errors"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/memory/vo"
)

var (
	ErrEmptyContent  = errors.New("memory content must not be empty")
	ErrMissingUserID = errors.New("userID is required")
)

// Memory is the aggregate root for a memory entry.
type Memory struct {
	id        vo.MemoryID
	userID    string
	content   string
	embedding []float32
	tags      []string
	createdAt time.Time
	updatedAt time.Time
}

// NewMemory creates a new Memory aggregate root.
// Returns ErrEmptyContent if content is empty, ErrMissingUserID if userID is empty.
func NewMemory(id vo.MemoryID, userID, content string, embedding []float32) (*Memory, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}
	if userID == "" {
		return nil, ErrMissingUserID
	}
	now := time.Now()
	return &Memory{
		id:        id,
		userID:    userID,
		content:   content,
		embedding: embedding,
		tags:      []string{},
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructMemory reconstructs a Memory from persistence without validation.
func ReconstructMemory(id vo.MemoryID, userID, content string, embedding []float32, tags []string, createdAt, updatedAt time.Time) *Memory {
	return &Memory{
		id:        id,
		userID:    userID,
		content:   content,
		embedding: embedding,
		tags:      tags,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// UpdateContent updates the memory content, embedding, and updatedAt timestamp.
func (m *Memory) UpdateContent(content string, embedding []float32) {
	m.content = content
	m.embedding = embedding
	m.updatedAt = time.Now()
}

// Getters
func (m *Memory) ID() vo.MemoryID      { return m.id }
func (m *Memory) UserID() string       { return m.userID }
func (m *Memory) Content() string      { return m.content }
func (m *Memory) Embedding() []float32 { return m.embedding }
func (m *Memory) Tags() []string       { return m.tags }
func (m *Memory) CreatedAt() time.Time { return m.createdAt }
func (m *Memory) UpdatedAt() time.Time { return m.updatedAt }
