package mapper

import (
	"encoding/json"
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/messaging/entity"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

// DomainToChannelPO converts a Channel entity to a ChannelPO.
func DomainToChannelPO(ch *entity.Channel) (*po.ChannelPO, error) {
	capsJSON, err := json.Marshal(ch.Capabilities())
	if err != nil {
		return nil, fmt.Errorf("marshal capabilities: %w", err)
	}
	return &po.ChannelPO{
		ChannelID:    ch.ID().Value(),
		ChannelType:  int(ch.ChannelType()),
		PluginID:     ch.PluginID(),
		Status:       int(ch.Status()),
		Capabilities: string(capsJSON),
	}, nil
}

// ChannelPOToDomain converts a ChannelPO to a Channel entity.
func ChannelPOToDomain(p *po.ChannelPO) (*entity.Channel, error) {
	channelID, err := vo.NewChannelID(p.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct channel_id: %w", err)
	}
	ch := entity.NewChannel(channelID, entity.ChannelType(p.ChannelType), p.PluginID)
	switch entity.ChannelStatus(p.Status) {
	case entity.ChannelStatusActive:
		ch.Activate()
	case entity.ChannelStatusError:
		ch.SetError()
	}
	return ch, nil
}
