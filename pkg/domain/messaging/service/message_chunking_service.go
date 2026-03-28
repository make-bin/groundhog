package service

import (
	"github.com/make-bin/groundhog/pkg/domain/messaging/entity"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// MessageChunkingService splits a message content string into chunks for a given channel type.
type MessageChunkingService interface {
	Chunk(content string, channelType entity.ChannelType) []vo.MessageChunk
}
