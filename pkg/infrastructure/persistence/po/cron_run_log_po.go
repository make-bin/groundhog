package po

import "time"

// CronRunLogPO is the persistence object for cron run log entries.
type CronRunLogPO struct {
	ID             int64  `gorm:"primaryKey;autoIncrement"`
	JobID          string `gorm:"type:varchar(100);not null;index"`
	Ts             int64  `gorm:"not null"`
	Action         string `gorm:"type:varchar(50);not null"`
	Status         string `gorm:"type:varchar(20);not null;index"`
	Error          string `gorm:"type:text"`
	Summary        string `gorm:"type:text"`
	SessionID      string `gorm:"type:varchar(200)"`
	RunAtMs        int64  `gorm:"not null"`
	DurationMs     int64
	NextRunAtMs    *int64
	Model          string `gorm:"type:varchar(100)"`
	Provider       string `gorm:"type:varchar(100)"`
	InputTokens    int
	OutputTokens   int
	TotalTokens    int
	DeliveryStatus string `gorm:"type:varchar(50)"`
	DeliveryError  string `gorm:"type:text"`
	CreatedAt      time.Time
}

// TableName returns the database table name for CronRunLogPO.
func (CronRunLogPO) TableName() string { return "cron_run_logs" }
