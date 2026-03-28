package service

import (
	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
)

// SystemPromptService is a domain service for building system prompts dynamically.
type SystemPromptService interface {
	// Build constructs a system prompt string from the agent configuration and skill list.
	Build(agent *entity.Agent, skills []string) string
}
