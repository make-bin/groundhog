package service

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

// ModelSelectionService is a domain service for selecting and falling back between models.
type ModelSelectionService interface {
	// Select returns an Agent configured for the given ModelConfig.
	Select(ctx context.Context, cfg vo.ModelConfig) (*entity.Agent, error)

	// Fallback returns an Agent for the next model in the fallback chain.
	Fallback(ctx context.Context, current string, chain []string) (*entity.Agent, error)
}
