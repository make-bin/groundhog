package adk

import (
	"context"
	"fmt"
	"sync"

	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// MCPServerCfg is the runtime config for one MCP server (mirrors config.MCPServerConfig).
type MCPServerCfg struct {
	Name    string
	Command string
	Args    []string
	Env     []string
}

// mcpTool wraps a single tool exposed by an MCP server.
type mcpTool struct {
	name   string
	client *MCPClient
}

func (t *mcpTool) Name() string { return t.name }

func (t *mcpTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	return t.client.CallTool(ctx, t.name, args)
}

// MCPManager starts and manages multiple MCP server subprocesses.
// It exposes all discovered tools as []Tool for use in RunnerAdapter.
type MCPManager struct {
	mu      sync.Mutex
	clients []*MCPClient
	tools   []Tool
	schemas []ToolSchema
	log     logger.Logger
}

// NewMCPManager creates an MCPManager. Call Start() to launch servers.
func NewMCPManager(log logger.Logger) *MCPManager {
	return &MCPManager{log: log}
}

// Start launches all configured MCP servers, discovers their tools, and stores them.
// Servers that fail to start are logged and skipped (best-effort).
func (m *MCPManager) Start(ctx context.Context, servers []MCPServerCfg) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, srv := range servers {
		client, err := NewMCPClient(ctx, srv.Name, srv.Command, srv.Args, srv.Env)
		if err != nil {
			m.log.Warn("mcp server failed to start", "name", srv.Name, "error", err)
			continue
		}

		defs, err := client.ListTools(ctx)
		if err != nil {
			m.log.Warn("mcp server tools/list failed", "name", srv.Name, "error", err)
			_ = client.Close()
			continue
		}

		m.clients = append(m.clients, client)
		for _, def := range defs {
			m.tools = append(m.tools, &mcpTool{name: fmt.Sprintf("%s__%s", srv.Name, def.Name), client: client})
			// Also register with bare name for single-server convenience
			m.tools = append(m.tools, &mcpTool{name: def.Name, client: client})

			schema := def.InputSchema
			if schema == nil {
				schema = map[string]any{"type": "object", "properties": map[string]any{}}
			}
			m.schemas = append(m.schemas, ToolSchema{
				Type: "function",
				Function: ToolSchemaFunc{
					Name:        def.Name,
					Description: def.Description,
					Parameters:  schema,
				},
			})
			m.log.Info("mcp tool registered", "server", srv.Name, "tool", def.Name)
		}
		m.log.Info("mcp server started", "name", srv.Name, "tools", len(defs))
	}
}

// Tools returns all tools discovered from MCP servers.
func (m *MCPManager) Tools() []Tool { m.mu.Lock(); defer m.mu.Unlock(); return m.tools }

// Schemas returns ToolSchemas for all MCP tools (for LLM function calling).
func (m *MCPManager) Schemas() []ToolSchema { m.mu.Lock(); defer m.mu.Unlock(); return m.schemas }

// Close shuts down all MCP server subprocesses.
func (m *MCPManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.clients {
		if err := c.Close(); err != nil {
			m.log.Warn("mcp client close error", "error", err)
		}
	}
	m.clients = nil
}
