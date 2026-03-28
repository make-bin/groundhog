package service

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

// ToolPolicyService is a domain service for evaluating tool execution policies.
type ToolPolicyService interface {
	// Evaluate returns the policy for executing the given tool with the provided arguments.
	Evaluate(ctx context.Context, tool entity.ToolDefinition, args map[string]any) (vo.ToolPolicy, error)
}
