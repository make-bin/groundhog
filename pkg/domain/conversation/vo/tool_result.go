// @AI_GENERATED
package vo

// ToolResult represents the result of a tool execution.
// It is immutable after creation.
type ToolResult struct {
	output  string
	isError bool
	errMsg  string
}

// NewToolResult creates a successful tool result with the given output.
func NewToolResult(output string) ToolResult {
	return ToolResult{output: output, isError: false, errMsg: ""}
}

// NewToolError creates an error tool result with the given error message.
func NewToolError(errMsg string) ToolResult {
	return ToolResult{output: "", isError: true, errMsg: errMsg}
}

// Output returns the tool output.
func (t ToolResult) Output() string { return t.output }

// IsError returns true if the tool execution resulted in an error.
func (t ToolResult) IsError() bool { return t.isError }

// ErrMsg returns the error message.
func (t ToolResult) ErrMsg() string { return t.errMsg }

// @AI_GENERATED: end
