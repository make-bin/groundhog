package po

import "time"

// ChannelPO is the persistence object for Channel entity.
type ChannelPO struct {
	ID           uint   `gorm:"primaryKey"`
	ChannelID    string `gorm:"uniqueIndex;type:varchar(100);not null"`
	ChannelType  int    `gorm:"not null"`
	PluginID     string `gorm:"type:varchar(100)"`
	Status       int    `gorm:"default:0"`
	Capabilities string `gorm:"type:text"` // JSON array
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (ChannelPO) TableName() string { return "channels" }
