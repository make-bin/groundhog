// @AI_GENERATED
package entity

import "github.com/make-bin/groundhog/pkg/domain/conversation/vo"

// Agent is an entity representing an AI agent configuration.
type Agent struct {
	id          vo.AgentID
	name        string
	description string
	modelConfig vo.ModelConfig
	tools       []ToolDefinition
	instruction string
	subAgents   []string
}

// NewAgent creates a new Agent with the given id and name.
func NewAgent(id vo.AgentID, name string) *Agent {
	return &Agent{
		id:   id,
		name: name,
	}
}

// ID returns the agent identifier.
func (a *Agent) ID() vo.AgentID { return a.id }

// Name returns the agent name.
func (a *Agent) Name() string { return a.name }

// Description returns the agent description.
func (a *Agent) Description() string { return a.description }

// ModelConfig returns the agent model configuration.
func (a *Agent) ModelConfig() vo.ModelConfig { return a.modelConfig }

// Tools returns a copy of the agent tool definitions.
func (a *Agent) Tools() []ToolDefinition {
	tools := make([]ToolDefinition, len(a.tools))
	copy(tools, a.tools)
	return tools
}

// Instruction returns the agent instruction.
func (a *Agent) Instruction() string { return a.instruction }

// SubAgents returns a copy of the sub-agent identifiers.
func (a *Agent) SubAgents() []string {
	subs := make([]string, len(a.subAgents))
	copy(subs, a.subAgents)
	return subs
}

// SetDescription sets the agent description.
func (a *Agent) SetDescription(description string) {
	a.description = description
}

// SetModelConfig sets the agent model configuration.
func (a *Agent) SetModelConfig(cfg vo.ModelConfig) {
	a.modelConfig = cfg
}

// SetInstruction sets the agent instruction.
func (a *Agent) SetInstruction(instruction string) {
	a.instruction = instruction
}

// SetTools sets the agent tool definitions.
func (a *Agent) SetTools(tools []ToolDefinition) {
	t := make([]ToolDefinition, len(tools))
	copy(t, tools)
	a.tools = t
}

// SetSubAgents sets the sub-agent identifiers.
func (a *Agent) SetSubAgents(subAgents []string) {
	s := make([]string, len(subAgents))
	copy(s, subAgents)
	a.subAgents = s
}

// @AI_GENERATED: end
