package service

import (
	"fmt"
	"strings"

	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

type commandDetectionServiceImpl struct{}

// NewCommandDetectionService creates a new CommandDetectionService implementation.
func NewCommandDetectionService() CommandDetectionService {
	return &commandDetectionServiceImpl{}
}

func (s *commandDetectionServiceImpl) Detect(content string) *vo.ParsedCommand {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "/") {
		return nil
	}
	parts := strings.Fields(trimmed[1:]) // strip leading /
	if len(parts) == 0 {
		return nil
	}
	name := parts[0]
	args := parts[1:]
	cmd := vo.NewParsedCommand(name, args)
	return &cmd
}

func (s *commandDetectionServiceImpl) Authorize(cmd *vo.ParsedCommand, principalID string) error {
	if cmd == nil {
		return fmt.Errorf("command is nil")
	}
	// Default: all principals are authorized for all commands
	// TODO: implement role-based authorization
	return nil
}
