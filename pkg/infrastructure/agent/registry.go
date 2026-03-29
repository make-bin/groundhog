// Package agent provides an in-memory registry of agent configurations loaded from config.
// It mirrors openclaw's agent-scope.ts resolveAgentConfig / resolveDefaultAgentId logic.
package agent

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/make-bin/groundhog/pkg/utils/config"
)

const DefaultAgentID = "main"

var validAgentIDRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

// Registry holds all agent configurations loaded from config, keyed by normalized ID.
type Registry struct {
	agents    map[string]*config.AgentEntryConfig
	defaults  config.AgentDefaultsConfig
	defaultID string
	ordered   []string // insertion order for determinism
}

// NewRegistry builds a Registry from the agents section of AppConfig.
// Mirrors openclaw resolveDefaultAgentId + listAgentEntries.
func NewRegistry(cfg *config.AppConfig) *Registry {
	r := &Registry{
		agents:   make(map[string]*config.AgentEntryConfig),
		defaults: cfg.Agents.Defaults,
	}

	for i := range cfg.Agents.List {
		entry := &cfg.Agents.List[i]
		id := NormalizeAgentID(entry.ID)
		if id == "" {
			continue
		}
		if _, exists := r.agents[id]; exists {
			continue // first definition wins (same as openclaw)
		}
		r.agents[id] = entry
		r.ordered = append(r.ordered, id)
	}

	r.defaultID = r.resolveDefaultID()
	return r
}

// resolveDefaultID returns the ID of the default agent.
// Priority: first entry with default:true → first entry → "main".
func (r *Registry) resolveDefaultID() string {
	for _, id := range r.ordered {
		if r.agents[id].Default {
			return id
		}
	}
	if len(r.ordered) > 0 {
		return r.ordered[0]
	}
	return DefaultAgentID
}

// DefaultID returns the resolved default agent ID.
func (r *Registry) DefaultID() string { return r.defaultID }

// Get returns the resolved agent config for the given ID, merging defaults.
// Returns an error if the agent is not found.
func (r *Registry) Get(agentID string) (*ResolvedAgent, error) {
	id := NormalizeAgentID(agentID)
	entry, ok := r.agents[id]
	if !ok {
		// If no agents configured at all, synthesize a "main" agent from global defaults.
		if len(r.agents) == 0 && id == DefaultAgentID {
			return r.syntheticMain(), nil
		}
		return nil, fmt.Errorf("agent %q not found", agentID)
	}
	return r.resolve(id, entry), nil
}

// GetDefault returns the resolved config for the default agent.
func (r *Registry) GetDefault() *ResolvedAgent {
	if len(r.agents) == 0 {
		return r.syntheticMain()
	}
	entry := r.agents[r.defaultID]
	return r.resolve(r.defaultID, entry)
}

// List returns all registered agents in insertion order.
func (r *Registry) List() []*ResolvedAgent {
	result := make([]*ResolvedAgent, 0, len(r.ordered))
	for _, id := range r.ordered {
		result = append(result, r.resolve(id, r.agents[id]))
	}
	return result
}

// resolve merges an entry with global defaults to produce a ResolvedAgent.
func (r *Registry) resolve(id string, entry *config.AgentEntryConfig) *ResolvedAgent {
	provider := firstNonEmpty(entry.Provider, r.defaults.Provider)
	model := firstNonEmpty(entry.Model, r.defaults.Model)
	systemPrompt := firstNonEmpty(entry.SystemPrompt, r.defaults.SystemPrompt)
	skills := entry.Skills
	if len(skills) == 0 {
		skills = r.defaults.Skills
	}
	workspace := firstNonEmpty(entry.Workspace, r.defaults.Workspace)

	name := entry.Name
	if name == "" {
		name = id
	}

	return &ResolvedAgent{
		ID:           id,
		Name:         name,
		Description:  entry.Description,
		Provider:     provider,
		Model:        model,
		SystemPrompt: systemPrompt,
		Skills:       skills,
		Workspace:    workspace,
		IsDefault:    id == r.defaultID,
	}
}

// syntheticMain returns a minimal agent built purely from global defaults.
// Used when no agents are configured in config.yaml.
func (r *Registry) syntheticMain() *ResolvedAgent {
	return &ResolvedAgent{
		ID:           DefaultAgentID,
		Name:         "Main",
		Provider:     r.defaults.Provider,
		Model:        r.defaults.Model,
		SystemPrompt: r.defaults.SystemPrompt,
		Skills:       r.defaults.Skills,
		Workspace:    r.defaults.Workspace,
		IsDefault:    true,
	}
}

// ResolvedAgent is the fully merged agent configuration ready for use.
type ResolvedAgent struct {
	ID           string
	Name         string
	Description  string
	Provider     string
	Model        string
	SystemPrompt string
	Skills       []string
	Workspace    string
	IsDefault    bool
}

// NormalizeAgentID lowercases and sanitizes an agent ID, mirroring openclaw normalizeAgentId.
func NormalizeAgentID(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return DefaultAgentID
	}
	lower := strings.ToLower(trimmed)
	if validAgentIDRe.MatchString(lower) {
		return lower
	}
	// Best-effort: collapse invalid chars to "-"
	invalidRe := regexp.MustCompile(`[^a-z0-9_-]+`)
	normalized := invalidRe.ReplaceAllString(lower, "-")
	normalized = strings.TrimLeft(normalized, "-")
	normalized = strings.TrimRight(normalized, "-")
	if len(normalized) > 64 {
		normalized = normalized[:64]
	}
	if normalized == "" {
		return DefaultAgentID
	}
	return normalized
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
