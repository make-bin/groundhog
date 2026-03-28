// @AI_GENERATED
package vo

import "strings"

// Prompt represents a template-based prompt with variable substitution.
// It is immutable after creation.
type Prompt struct {
	template string
	vars     map[string]string
}

// NewPrompt creates a new Prompt with the given template and variables.
func NewPrompt(template string, vars map[string]string) Prompt {
	// Copy vars to ensure immutability.
	v := make(map[string]string, len(vars))
	for k, val := range vars {
		v[k] = val
	}
	return Prompt{template: template, vars: v}
}

// Template returns the prompt template.
func (p Prompt) Template() string { return p.template }

// Vars returns a copy of the prompt variables.
func (p Prompt) Vars() map[string]string {
	v := make(map[string]string, len(p.vars))
	for k, val := range p.vars {
		v[k] = val
	}
	return v
}

// Render replaces {{key}} placeholders in the template with values from vars.
func (p Prompt) Render() string {
	result := p.template
	for k, v := range p.vars {
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result
}

// @AI_GENERATED: end
