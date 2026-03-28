// @AI_GENERATED
package vo

// SessionState represents the state of an agent session.
type SessionState int

const (
	SessionStateActive SessionState = iota
	SessionStatePaused
	SessionStateArchived
	SessionStateCompacted
)

// String returns the string representation of the session state.
func (s SessionState) String() string {
	switch s {
	case SessionStateActive:
		return "Active"
	case SessionStatePaused:
		return "Paused"
	case SessionStateArchived:
		return "Archived"
	case SessionStateCompacted:
		return "Compacted"
	default:
		return "unknown"
	}
}

// @AI_GENERATED: end
