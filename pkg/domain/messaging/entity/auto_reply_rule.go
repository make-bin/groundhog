// @AI_GENERATED
package entity

// AutoReplyRule represents a rule that triggers an automatic reply based on message content.
type AutoReplyRule struct {
	id        string
	trigger   string
	authLevel int
	handler   string
}

// NewAutoReplyRule creates a new AutoReplyRule.
func NewAutoReplyRule(id, trigger string, authLevel int, handler string) *AutoReplyRule {
	return &AutoReplyRule{
		id:        id,
		trigger:   trigger,
		authLevel: authLevel,
		handler:   handler,
	}
}

// ID returns the rule identifier.
func (a *AutoReplyRule) ID() string { return a.id }

// Trigger returns the trigger pattern for this rule.
func (a *AutoReplyRule) Trigger() string { return a.trigger }

// AuthLevel returns the minimum authorization level required to trigger this rule.
func (a *AutoReplyRule) AuthLevel() int { return a.authLevel }

// Handler returns the handler identifier for this rule.
func (a *AutoReplyRule) Handler() string { return a.handler }

// @AI_GENERATED: end
