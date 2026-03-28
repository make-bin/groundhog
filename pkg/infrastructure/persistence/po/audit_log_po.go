// @AI_GENERATED
package po

import "time"

// AuditLogPO is the persistence object for audit log entries.
type AuditLogPO struct {
	ID           uint      `gorm:"primaryKey"`
	Action       string    `gorm:"type:varchar(100);not null;index"`
	PrincipalID  string    `gorm:"type:varchar(100);index"`
	ResourceType string    `gorm:"type:varchar(100)"`
	ResourceID   string    `gorm:"type:varchar(100)"`
	Details      string    `gorm:"type:text"`
	SourceIP     string    `gorm:"type:varchar(50)"`
	CreatedAt    time.Time `gorm:"index"`
}

func (AuditLogPO) TableName() string { return "audit_logs" }

// @AI_GENERATED: end
