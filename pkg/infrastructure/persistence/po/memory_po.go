package po

import (
	"time"

	"gorm.io/gorm"
)

// MemoryPO is the persistence object for the Memory aggregate.
// Embedding is stored as TEXT to remain compatible with databases that
// do not have the pgvector extension installed. When pgvector is available
// the migration will ALTER the column to VECTOR(1024) automatically.
type MemoryPO struct {
	ID        uint   `gorm:"primaryKey"`
	MemoryID  string `gorm:"uniqueIndex;type:varchar(100);not null"`
	UserID    string `gorm:"type:varchar(100);not null;index"`
	Content   string `gorm:"type:text;not null"`
	Embedding string `gorm:"type:text"` // JSON float array or pgvector literal
	Tags      string `gorm:"type:text"` // JSON array
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the database table name for MemoryPO.
func (MemoryPO) TableName() string { return "memories" }
