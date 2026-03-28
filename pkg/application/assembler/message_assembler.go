package assembler

import (
	"github.com/make-bin/groundhog/pkg/application/dto"
	"github.com/make-bin/groundhog/pkg/domain/messaging/aggregate/inbound_message"
	"github.com/make-bin/groundhog/pkg/domain/messaging/entity"
)

// ToMessageDTO converts an InboundMessage aggregate to a MessageDTO.
func ToMessageDTO(msg *inbound_message.InboundMessage) *dto.MessageDTO {
	return &dto.MessageDTO{
		ID:         msg.ID().Value(),
		ChannelID:  msg.ChannelID().Value(),
		AccountID:  msg.AccountID().Value(),
		Content:    msg.Content().Text(),
		Status:     msg.Status().String(),
		RoutedTo:   msg.RoutedTo(),
		ReceivedAt: msg.ReceivedAt(),
	}
}

// ToChannelDTO converts a Channel entity to a ChannelDTO.
func ToChannelDTO(ch *entity.Channel) *dto.ChannelDTO {
	return &dto.ChannelDTO{
		ID:           ch.ID().Value(),
		ChannelType:  ch.ChannelType().String(),
		PluginID:     ch.PluginID(),
		Status:       ch.Status().String(),
		Capabilities: ch.Capabilities(),
	}
}

// ToChannelDTOList converts a slice of Channel entities to a slice of ChannelDTOs.
func ToChannelDTOList(channels []*entity.Channel) []*dto.ChannelDTO {
	result := make([]*dto.ChannelDTO, 0, len(channels))
	for _, ch := range channels {
		result = append(result, ToChannelDTO(ch))
	}
	return result
}
