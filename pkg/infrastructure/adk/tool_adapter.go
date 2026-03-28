// @AI_GENERATED
package adk

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
)

// Tool is the internal interface for an executable tool.
type Tool interface {
	Name() string
	Execute(ctx context.Context, args map[string]any) (string, error)
}

// BeforeToolCallback is called before a tool executes. Return an error to deny execution.
type BeforeToolCallback func(ctx context.Context, toolName string, args map[string]any) error

// ApprovalNotifier is called when a tool is blocked waiting for approval.
type ApprovalNotifier func(pa interface{})

// ToolAdapter converts domain ToolDefinitions to internal Tool instances.
type ToolAdapter struct {
	builtins         map[string]Tool
	mcpTools         map[string]Tool // tools from MCP servers
	mcpSchemas       []ToolSchema    // schemas for MCP tools (for LLM function calling)
	beforeCallback   BeforeToolCallback
	approvalNotifier ApprovalNotifier
}

// NewToolAdapter creates a ToolAdapter with all built-in tools registered.
func NewToolAdapter() *ToolAdapter {
	a := &ToolAdapter{
		builtins: make(map[string]Tool),
		mcpTools: make(map[string]Tool),
	}
	for _, t := range []Tool{
		&bashExecTool{},
		&fileReadTool{},
		&fileWriteTool{},
		&fileListTool{},
		&webSearchTool{},
	} {
		a.builtins[t.Name()] = t
	}
	return a
}

// RegisterMCPTools registers tools discovered from MCP servers.
// Call this after MCPManager.Start() during server initialization.
func (a *ToolAdapter) RegisterMCPTools(mgr *MCPManager) {
	for _, t := range mgr.Tools() {
		a.mcpTools[t.Name()] = t
	}
	a.mcpSchemas = append(a.mcpSchemas, mgr.Schemas()...)
}

// MCPSchemas returns ToolSchemas for all registered MCP tools.
func (a *ToolAdapter) MCPSchemas() []ToolSchema {
	return a.mcpSchemas
}

// SetBeforeToolCallback sets the callback invoked before each tool execution.
func (a *ToolAdapter) SetBeforeToolCallback(cb BeforeToolCallback) {
	a.beforeCallback = cb
}

// SetApprovalNotifier sets the notifier called when a tool is pending approval.
func (a *ToolAdapter) SetApprovalNotifier(fn ApprovalNotifier) {
	a.approvalNotifier = fn
}

// ToADKTools converts domain ToolDefinitions to Tool instances.
// MCP tools are always included regardless of session tool definitions.
func (a *ToolAdapter) ToADKTools(tools []entity.ToolDefinition) []Tool {
	result := make([]Tool, 0, len(tools)+len(a.mcpTools))
	for _, td := range tools {
		if t, ok := a.builtins[td.Name()]; ok {
			result = append(result, t)
		}
	}
	// Append MCP tools (deduplicate by name)
	for _, t := range a.mcpTools {
		result = append(result, t)
	}
	return result
}

// ToToolSchemas converts domain ToolDefinitions to OpenAI function calling format ToolSchemas.
// MCP tool schemas are always appended.
func (a *ToolAdapter) ToToolSchemas(tools []entity.ToolDefinition) []ToolSchema {
	schemas := make([]ToolSchema, 0, len(tools)+len(a.mcpSchemas))
	for _, td := range tools {
		params := td.Schema()
		if len(params) == 0 {
			params = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		schemas = append(schemas, ToolSchema{
			Type: "function",
			Function: ToolSchemaFunc{
				Name:        td.Name(),
				Description: td.Description(),
				Parameters:  params,
			},
		})
	}
	schemas = append(schemas, a.mcpSchemas...)
	return schemas
}

// --- bash_exec ---

const bashTimeout = 30 * time.Second
const bashMaxOutput = 10 * 1024 // 10KB

type bashExecTool struct{}

func (b *bashExecTool) Name() string { return "bash_exec" }

func (b *bashExecTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	cmd, _ := args["command"].(string)
	if cmd == "" {
		return "", fmt.Errorf("bash_exec: 'command' argument is required")
	}
	ctx, cancel := context.WithTimeout(ctx, bashTimeout)
	defer cancel()
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	_ = c.Run()
	result := out.String()
	if len(result) > bashMaxOutput {
		result = result[:bashMaxOutput] + "\n[output truncated]"
	}
	return result, nil
}

// --- file_read ---

type fileReadTool struct{}

func (f *fileReadTool) Name() string { return "file_read" }

func (f *fileReadTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("file_read: 'path' argument is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("file_read: %w", err)
	}
	return string(data), nil
}

// --- file_write ---

type fileWriteTool struct{}

func (f *fileWriteTool) Name() string { return "file_write" }

func (f *fileWriteTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	content, _ := args["content"].(string)
	if path == "" {
		return "", fmt.Errorf("file_write: 'path' argument is required")
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("file_write: %w", err)
	}
	return fmt.Sprintf("wrote %d bytes to %s", len(content), path), nil
}

// --- file_list ---

type fileListTool struct{}

func (f *fileListTool) Name() string { return "file_list" }

func (f *fileListTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	dir, _ := args["dir"].(string)
	if dir == "" {
		dir = "."
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("file_list: %w", err)
	}
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(filepath.Join(dir, e.Name()))
		if e.IsDir() {
			sb.WriteString("/")
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

// --- web_search ---

type webSearchTool struct{}

func (w *webSearchTool) Name() string { return "web_search" }

func (w *webSearchTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	query, _ := args["query"].(string)
	if query == "" {
		return "", fmt.Errorf("web_search: 'query' argument is required")
	}
	// Stub: real implementation would call a search API
	return fmt.Sprintf("[web_search stub] query: %s", query), nil
}

// @AI_GENERATED: end
