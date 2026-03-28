package po

import "time"

// MessagePO is the persistence object for InboundMessage aggregate.
type MessagePO struct {
	ID         uint   `gorm:"primaryKey"`
	MessageID  string `gorm:"uniqueIndex;type:varchar(100);not null"`
	ChannelID  string `gorm:"type:varchar(100);not null;index"`
	AccountID  string `gorm:"type:varchar(100);not null;index"`
	Content    string `gorm:"type:text"`
	Status     int    `gorm:"default:0"`
	RoutedTo   string `gorm:"type:varchar(100)"`
	ReceivedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (MessagePO) TableName() string { return "messages" }
