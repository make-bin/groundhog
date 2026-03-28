// @AI_GENERATED
package entity

import "github.com/make-bin/groundhog/pkg/domain/conversation/vo"

// ToolDefinition is an entity representing a tool available to an agent.
type ToolDefinition struct {
	name        string
	description string
	category    vo.ToolCategory
	schema      map[string]any
	policy      vo.ToolPolicy
}

// NewToolDefinition creates a new ToolDefinition with the given parameters.
func NewToolDefinition(name, description string, category vo.ToolCategory, schema map[string]any, policy vo.ToolPolicy) *ToolDefinition {
	s := make(map[string]any, len(schema))
	for k, v := range schema {
		s[k] = v
	}
	return &ToolDefinition{
		name:        name,
		description: description,
		category:    category,
		schema:      s,
		policy:      policy,
	}
}

// Name returns the tool name.
func (t *ToolDefinition) Name() string { return t.name }

// Description returns the tool description.
func (t *ToolDefinition) Description() string { return t.description }

// Category returns the tool category.
func (t *ToolDefinition) Category() vo.ToolCategory { return t.category }

// Schema returns a copy of the tool schema.
func (t *ToolDefinition) Schema() map[string]any {
	s := make(map[string]any, len(t.schema))
	for k, v := range t.schema {
		s[k] = v
	}
	return s
}

// Policy returns the tool execution policy.
func (t *ToolDefinition) Policy() vo.ToolPolicy { return t.policy }

// @AI_GENERATED: end
