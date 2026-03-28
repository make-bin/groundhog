// @AI_GENERATED
package entity

import "time"

// AuditEntry is an entity representing an audit log entry.
type AuditEntry struct {
	action    string
	sourceIP  string
	timestamp time.Time
}

// NewAuditEntry creates a new AuditEntry with the given action and source IP.
// The timestamp is set to the current time.
func NewAuditEntry(action, sourceIP string) *AuditEntry {
	return &AuditEntry{
		action:    action,
		sourceIP:  sourceIP,
		timestamp: time.Now(),
	}
}

// Action returns the audit action.
func (a *AuditEntry) Action() string { return a.action }

// SourceIP returns the source IP address.
func (a *AuditEntry) SourceIP() string { return a.sourceIP }

// Timestamp returns the time the audit entry was created.
func (a *AuditEntry) Timestamp() time.Time { return a.timestamp }

// @AI_GENERATED: end
