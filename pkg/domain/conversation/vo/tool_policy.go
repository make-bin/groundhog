// @AI_GENERATED
package vo

// ToolPolicy represents the execution policy for a tool.
type ToolPolicy int

const (
	ToolPolicyAuto ToolPolicy = iota
	ToolPolicyConfirm
	ToolPolicyDeny
)

// String returns the string representation of the tool policy.
func (t ToolPolicy) String() string {
	switch t {
	case ToolPolicyAuto:
		return "auto"
	case ToolPolicyConfirm:
		return "confirm"
	case ToolPolicyDeny:
		return "deny"
	default:
		return "unknown"
	}
}

// @AI_GENERATED: end
