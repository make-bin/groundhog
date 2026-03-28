// @AI_GENERATED
package vo

// ToolCategory represents the category of a tool.
type ToolCategory int

const (
	ToolCategoryBash ToolCategory = iota
	ToolCategoryFile
	ToolCategoryWeb
	ToolCategoryBrowser
	ToolCategoryMCP
	ToolCategoryCustom
)

// String returns the string representation of the tool category.
func (t ToolCategory) String() string {
	switch t {
	case ToolCategoryBash:
		return "bash"
	case ToolCategoryFile:
		return "file"
	case ToolCategoryWeb:
		return "web"
	case ToolCategoryBrowser:
		return "browser"
	case ToolCategoryMCP:
		return "mcp"
	case ToolCategoryCustom:
		return "custom"
	default:
		return "unknown"
	}
}

// @AI_GENERATED: end
