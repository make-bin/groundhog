package adk

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

// ── JSON-RPC 2.0 wire types ───────────────────────────────────────────────────

type jsonrpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *jsonrpcError) Error() string { return fmt.Sprintf("mcp rpc error %d: %s", e.Code, e.Message) }

// ── MCP protocol types ────────────────────────────────────────────────────────

type mcpInitializeParams struct {
	ProtocolVersion string         `json:"protocolVersion"`
	ClientInfo      mcpClientInfo  `json:"clientInfo"`
	Capabilities    map[string]any `json:"capabilities"`
}

type mcpClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type mcpToolsListResult struct {
	Tools []MCPToolDef `json:"tools"`
}

// MCPToolDef is the tool definition returned by tools/list.
type MCPToolDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type mcpToolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type mcpToolCallResult struct {
	Content []mcpContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

type mcpContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ── MCPClient ─────────────────────────────────────────────────────────────────

// MCPClient manages a single MCP server subprocess and communicates via JSON-RPC over stdio.
type MCPClient struct {
	name    string
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	scanner *bufio.Scanner

	mu      sync.Mutex
	pending map[int64]chan jsonrpcResponse
	nextID  atomic.Int64

	done chan struct{}
}

// NewMCPClient starts the MCP server subprocess and performs the initialize handshake.
// env is a list of "KEY=VALUE" strings to add to the subprocess environment.
func NewMCPClient(ctx context.Context, name, command string, args, env []string) (*MCPClient, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = append(os.Environ(), env...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("mcp %s: stdin pipe: %w", name, err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("mcp %s: stdout pipe: %w", name, err)
	}
	// Discard stderr to avoid blocking
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("mcp %s: start process: %w", name, err)
	}

	c := &MCPClient{
		name:    name,
		cmd:     cmd,
		stdin:   stdin,
		scanner: bufio.NewScanner(stdout),
		pending: make(map[int64]chan jsonrpcResponse),
		done:    make(chan struct{}),
	}
	// Set a generous scanner buffer for large responses
	c.scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	go c.readLoop()

	// MCP initialize handshake
	initParams := mcpInitializeParams{
		ProtocolVersion: "2024-11-05",
		ClientInfo:      mcpClientInfo{Name: "groundhog", Version: "1.0.0"},
		Capabilities:    map[string]any{},
	}
	var initResult json.RawMessage
	if err := c.call(ctx, "initialize", initParams, &initResult); err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("mcp %s: initialize: %w", name, err)
	}
	// Send initialized notification (no response expected)
	_ = c.notify("notifications/initialized", nil)

	return c, nil
}

// ListTools calls tools/list and returns all available tool definitions.
func (c *MCPClient) ListTools(ctx context.Context) ([]MCPToolDef, error) {
	var result mcpToolsListResult
	if err := c.call(ctx, "tools/list", nil, &result); err != nil {
		return nil, fmt.Errorf("mcp %s: tools/list: %w", c.name, err)
	}
	return result.Tools, nil
}

// CallTool calls tools/call and returns the text output.
func (c *MCPClient) CallTool(ctx context.Context, toolName string, args map[string]any) (string, error) {
	params := mcpToolCallParams{Name: toolName, Arguments: args}
	var result mcpToolCallResult
	if err := c.call(ctx, "tools/call", params, &result); err != nil {
		return "", fmt.Errorf("mcp %s: tools/call %s: %w", c.name, toolName, err)
	}
	if result.IsError {
		for _, c := range result.Content {
			if c.Type == "text" && c.Text != "" {
				return "", fmt.Errorf("mcp tool error: %s", c.Text)
			}
		}
		return "", fmt.Errorf("mcp tool %s returned error", toolName)
	}
	var out string
	for _, content := range result.Content {
		if content.Type == "text" {
			out += content.Text
		}
	}
	return out, nil
}

// Close terminates the subprocess and cleans up resources.
func (c *MCPClient) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	_ = c.stdin.Close()
	return c.cmd.Process.Kill()
}

// ── internal ──────────────────────────────────────────────────────────────────

func (c *MCPClient) call(ctx context.Context, method string, params any, result any) error {
	id := c.nextID.Add(1)
	ch := make(chan jsonrpcResponse, 1)

	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	req := jsonrpcRequest{JSONRPC: "2.0", ID: id, Method: method, Params: params}
	data, err := json.Marshal(req)
	if err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return err
	}
	data = append(data, '\n')

	c.mu.Lock()
	_, writeErr := c.stdin.Write(data)
	c.mu.Unlock()
	if writeErr != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return fmt.Errorf("write request: %w", writeErr)
	}

	select {
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return ctx.Err()
	case resp := <-ch:
		if resp.Error != nil {
			return resp.Error
		}
		if result != nil && resp.Result != nil {
			return json.Unmarshal(resp.Result, result)
		}
		return nil
	case <-time.After(30 * time.Second):
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return fmt.Errorf("timeout waiting for response to %s", method)
	}
}

func (c *MCPClient) notify(method string, params any) error {
	// Notifications have no ID
	type notification struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  any    `json:"params,omitempty"`
	}
	data, err := json.Marshal(notification{JSONRPC: "2.0", Method: method, Params: params})
	if err != nil {
		return err
	}
	data = append(data, '\n')
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err = c.stdin.Write(data)
	return err
}

func (c *MCPClient) readLoop() {
	for c.scanner.Scan() {
		line := c.scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var resp jsonrpcResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			continue // skip malformed lines (e.g. server log output)
		}
		c.mu.Lock()
		ch, ok := c.pending[resp.ID]
		if ok {
			delete(c.pending, resp.ID)
		}
		c.mu.Unlock()
		if ok {
			ch <- resp
		}
	}
}
