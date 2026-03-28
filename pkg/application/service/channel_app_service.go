// @AI_GENERATED
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/application/dto"
	"github.com/make-bin/groundhog/pkg/domain/messaging/entity"
	"github.com/make-bin/groundhog/pkg/domain/messaging/repository"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// ChannelAppService defines the application service interface for channel management.
type ChannelAppService interface {
	ListChannels(ctx context.Context) ([]*dto.ChannelDTO, error)
	CreateChannel(ctx context.Context, req *dto.CreateChannelRequest) (*dto.ChannelDTO, error)
	DeleteChannel(ctx context.Context, id vo.ChannelID) error
	GetChannelStatus(ctx context.Context, id vo.ChannelID) (*dto.ChannelDTO, error)
}

type channelAppService struct {
	ChannelRepo repository.ChannelRepository `inject:""`
}

// NewChannelAppService creates a new ChannelAppService.
func NewChannelAppService() ChannelAppService {
	return &channelAppService{}
}

func (s *channelAppService) ListChannels(ctx context.Context) ([]*dto.ChannelDTO, error) {
	channels, err := s.ChannelRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.ChannelDTO, 0, len(channels))
	for _, ch := range channels {
		result = append(result, toChannelDTO(ch))
	}
	return result, nil
}

func (s *channelAppService) CreateChannel(ctx context.Context, req *dto.CreateChannelRequest) (*dto.ChannelDTO, error) {
	channelID, err := vo.NewChannelID(fmt.Sprintf("ch-%d", time.Now().UnixNano()))
	if err != nil {
		return nil, err
	}
	if req.ChannelID != "" {
		channelID, err = vo.NewChannelID(req.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	chType := channelTypeFromString(req.ChannelType)
	ch := entity.NewChannel(channelID, chType, req.PluginID)

	if err := s.ChannelRepo.Save(ctx, ch); err != nil {
		return nil, err
	}
	return toChannelDTO(ch), nil
}

func (s *channelAppService) DeleteChannel(ctx context.Context, id vo.ChannelID) error {
	return s.ChannelRepo.Delete(ctx, id)
}

func (s *channelAppService) GetChannelStatus(ctx context.Context, id vo.ChannelID) (*dto.ChannelDTO, error) {
	ch, err := s.ChannelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toChannelDTO(ch), nil
}

func toChannelDTO(ch *entity.Channel) *dto.ChannelDTO {
	return &dto.ChannelDTO{
		ID:           ch.ID().Value(),
		ChannelType:  ch.ChannelType().String(),
		PluginID:     ch.PluginID(),
		Status:       ch.Status().String(),
		Capabilities: ch.Capabilities(),
	}
}

func channelTypeFromString(s string) entity.ChannelType {
	switch s {
	case "discord", "Discord":
		return entity.ChannelTypeDiscord
	case "whatsapp", "WhatsApp":
		return entity.ChannelTypeWhatsApp
	case "slack", "Slack":
		return entity.ChannelTypeSlack
	default:
		return entity.ChannelTypeTelegram
	}
}

// @AI_GENERATED: end
