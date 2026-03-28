// @AI_GENERATED
package service

import (
	"fmt"
	"strings"

	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

type systemPromptServiceImpl struct{}

// NewSystemPromptService creates a new SystemPromptService implementation.
func NewSystemPromptService() SystemPromptService {
	return &systemPromptServiceImpl{}
}

// Build constructs a system prompt from the agent's instruction, tools, and skills.
// It renders the instruction as a Prompt template, then appends tool descriptions
// and skill identifiers.
func (s *systemPromptServiceImpl) Build(agent *entity.Agent, skills []string) string {
	if agent == nil {
		return ""
	}

	var sb strings.Builder

	// Render the agent instruction as a prompt template with no extra vars
	instruction := vo.NewPrompt(agent.Instruction(), nil).Render()
	if instruction != "" {
		sb.WriteString(instruction)
		sb.WriteString("\n\n")
	}

	// Append tool descriptions
	tools := agent.Tools()
	if len(tools) > 0 {
		sb.WriteString("## Available Tools\n\n")
		for _, t := range tools {
			sb.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", t.Name(), t.Category(), t.Description()))
		}
		sb.WriteString("\n")
	}

	// Append skills
	if len(skills) > 0 {
		sb.WriteString("## Active Skills\n\n")
		for _, skill := range skills {
			sb.WriteString(fmt.Sprintf("- %s\n", skill))
		}
		sb.WriteString("\n")
	}

	return strings.TrimSpace(sb.String())
}

// @AI_GENERATED: end
