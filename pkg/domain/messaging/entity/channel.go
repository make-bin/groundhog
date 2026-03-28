// @AI_GENERATED
package entity

import "github.com/make-bin/groundhog/pkg/domain/messaging/vo"

// ChannelType represents the type of messaging channel.
type ChannelType int

const (
	ChannelTypeTelegram ChannelType = iota
	ChannelTypeDiscord
	ChannelTypeWhatsApp
	ChannelTypeSlack
)

// String returns the string representation of the channel type.
func (t ChannelType) String() string {
	switch t {
	case ChannelTypeTelegram:
		return "Telegram"
	case ChannelTypeDiscord:
		return "Discord"
	case ChannelTypeWhatsApp:
		return "WhatsApp"
	case ChannelTypeSlack:
		return "Slack"
	default:
		return "unknown"
	}
}

// ChannelStatus represents the operational status of a channel.
type ChannelStatus int

const (
	ChannelStatusInactive ChannelStatus = iota
	ChannelStatusActive
	ChannelStatusError
)

// String returns the string representation of the channel status.
func (s ChannelStatus) String() string {
	switch s {
	case ChannelStatusInactive:
		return "Inactive"
	case ChannelStatusActive:
		return "Active"
	case ChannelStatusError:
		return "Error"
	default:
		return "unknown"
	}
}

// Channel represents a messaging channel entity.
type Channel struct {
	id           vo.ChannelID
	channelType  ChannelType
	pluginID     string
	status       ChannelStatus
	capabilities []string
}

// NewChannel creates a new Channel in Inactive status.
func NewChannel(id vo.ChannelID, channelType ChannelType, pluginID string) *Channel {
	return &Channel{
		id:           id,
		channelType:  channelType,
		pluginID:     pluginID,
		status:       ChannelStatusInactive,
		capabilities: []string{},
	}
}

// ID returns the channel identifier.
func (c *Channel) ID() vo.ChannelID { return c.id }

// ChannelType returns the channel type.
func (c *Channel) ChannelType() ChannelType { return c.channelType }

// PluginID returns the plugin identifier managing this channel.
func (c *Channel) PluginID() string { return c.pluginID }

// Status returns the current channel status.
func (c *Channel) Status() ChannelStatus { return c.status }

// Capabilities returns a copy of the channel capabilities.
func (c *Channel) Capabilities() []string {
	result := make([]string, len(c.capabilities))
	copy(result, c.capabilities)
	return result
}

// Activate sets the channel status to Active.
func (c *Channel) Activate() {
	c.status = ChannelStatusActive
}

// Deactivate sets the channel status to Inactive.
func (c *Channel) Deactivate() {
	c.status = ChannelStatusInactive
}

// SetError sets the channel status to Error.
func (c *Channel) SetError() {
	c.status = ChannelStatusError
}

// @AI_GENERATED: end
