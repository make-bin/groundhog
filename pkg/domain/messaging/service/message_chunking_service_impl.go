package service

import (
	"github.com/make-bin/groundhog/pkg/domain/messaging/entity"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// channelLimits maps ChannelType to its character limit.
var channelLimits = map[entity.ChannelType]int{
	entity.ChannelTypeTelegram: 4096,
	entity.ChannelTypeDiscord:  2000,
	entity.ChannelTypeWhatsApp: 4096,
	entity.ChannelTypeSlack:    40000,
}

const defaultLimit = 4096

type messageChunkingServiceImpl struct{}

// NewMessageChunkingService creates a new MessageChunkingService implementation.
func NewMessageChunkingService() MessageChunkingService {
	return &messageChunkingServiceImpl{}
}

func (s *messageChunkingServiceImpl) Chunk(content string, channelType entity.ChannelType) []vo.MessageChunk {
	limit, ok := channelLimits[channelType]
	if !ok {
		limit = defaultLimit
	}

	if len(content) <= limit {
		return []vo.MessageChunk{vo.NewMessageChunk(0, content, true)}
	}

	var chunks []vo.MessageChunk
	runes := []rune(content)
	total := len(runes)
	idx := 0
	for start := 0; start < total; start += limit {
		end := start + limit
		if end > total {
			end = total
		}
		isFinal := end >= total
		chunks = append(chunks, vo.NewMessageChunk(idx, string(runes[start:end]), isFinal))
		idx++
	}
	return chunks
}
