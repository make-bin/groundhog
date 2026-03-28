package po

import (
	"time"

	"gorm.io/gorm"
)

// CronJobPO is the persistence object for the CronJob aggregate.
type CronJobPO struct {
	ID             uint   `gorm:"primaryKey"`
	JobID          string `gorm:"uniqueIndex;type:varchar(100);not null"`
	AgentId        string `gorm:"type:varchar(100);index"`
	SessionKey     string `gorm:"type:varchar(200)"`
	Name           string `gorm:"type:varchar(200);uniqueIndex;not null"`
	Description    string `gorm:"type:text"`
	Enabled        bool   `gorm:"default:true;index"`
	DeleteAfterRun *bool
	Schedule       string `gorm:"type:jsonb;not null"`
	SessionTarget  string `gorm:"type:varchar(200);not null"`
	WakeMode       string `gorm:"type:varchar(50);not null"`
	Payload        string `gorm:"type:jsonb;not null"`
	Delivery       string `gorm:"type:jsonb"`
	FailureAlert   string `gorm:"type:jsonb"`
	State          string `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// TableName returns the database table name for CronJobPO.
func (CronJobPO) TableName() string { return "cron_jobs" }
