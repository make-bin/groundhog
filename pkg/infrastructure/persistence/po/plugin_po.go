package po

import "time"

// PluginPO is the persistence object for PluginInstance aggregate.
type PluginPO struct {
	ID           uint   `gorm:"primaryKey"`
	PluginID     string `gorm:"uniqueIndex;type:varchar(100);not null"`
	Name         string `gorm:"type:varchar(200);not null"`
	Version      string `gorm:"type:varchar(50)"`
	PluginType   string `gorm:"type:varchar(50)"`
	EntryPoint   string `gorm:"type:varchar(500)"`
	Status       int    `gorm:"default:0"`
	Capabilities string `gorm:"type:text"` // JSON array
	StartedAt    *time.Time
	RestartCount int `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (PluginPO) TableName() string { return "plugins" }
