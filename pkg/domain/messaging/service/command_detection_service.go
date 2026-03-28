package service

import (
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// CommandDetectionService detects and authorizes chat commands.
type CommandDetectionService interface {
	Detect(content string) *vo.ParsedCommand
	Authorize(cmd *vo.ParsedCommand, principalID string) error
}
