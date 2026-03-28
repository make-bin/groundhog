// @AI_GENERATED
package po

import (
	"time"

	"gorm.io/gorm"
)

// SessionPO is the persistence object for AgentSession aggregate.
type SessionPO struct {
	ID          uint   `gorm:"primaryKey"`
	SessionID   string `gorm:"uniqueIndex;type:varchar(100);not null"`
	AgentID     string `gorm:"type:varchar(100);not null"`
	UserID      string `gorm:"type:varchar(100);not null;index"`
	ActiveModel string `gorm:"type:text"` // JSON
	State       int    `gorm:"default:0"`
	Metadata    string `gorm:"type:text"` // JSON
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Turns       []TurnPO       `gorm:"foreignKey:SessionPOID"`
}

// TableName returns the database table name for SessionPO.
func (SessionPO) TableName() string { return "sessions" }

// TurnPO is the persistence object for Turn entity.
type TurnPO struct {
	ID          uint   `gorm:"primaryKey"`
	SessionPOID uint   `gorm:"index"`
	TurnID      string `gorm:"type:varchar(100)"`
	UserInput   string `gorm:"type:text"`
	Response    string `gorm:"type:text"`
	ModelUsed   string `gorm:"type:varchar(100)"`
	TokenUsage  string `gorm:"type:text"` // JSON
	ToolCalls   string `gorm:"type:text"` // JSON
	StartedAt   time.Time
	CompletedAt time.Time
}

// TableName returns the database table name for TurnPO.
func (TurnPO) TableName() string { return "turns" }

// @AI_GENERATED: end
