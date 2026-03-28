package grpc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	channelv1 "github.com/make-bin/groundhog/proto/gen/channel/v1"
)

// NodeBridge manages a Node.js subprocess that hosts existing JS channel plugins.
type NodeBridge struct {
	extensionPath string
	cmd           *exec.Cmd
	msgCh         chan *channelv1.InboundMessageProto
	mu            sync.Mutex
}

// NewNodeBridge creates a new NodeBridge for the given extensions directory.
func NewNodeBridge(extensionPath string) (*NodeBridge, error) {
	if _, err := os.Stat(extensionPath); err != nil {
		return nil, fmt.Errorf("node_bridge: extensions path %q not found: %w", extensionPath, err)
	}
	return &NodeBridge{
		extensionPath: extensionPath,
		msgCh:         make(chan *channelv1.InboundMessageProto, 100),
	}, nil
}

// Start launches the Node.js subprocess with the bridge server script.
func (b *NodeBridge) Start(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// The Node.js bridge script is expected at extensions/bridge/server.js
	scriptPath := b.extensionPath + "/bridge/server.js"
	cmd := exec.CommandContext(ctx, "node", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"EXTENSIONS_PATH="+b.extensionPath,
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("node_bridge: start node process: %w", err)
	}
	b.cmd = cmd

	go func() {
		_ = cmd.Wait()
	}()

	return nil
}

// Stop terminates the Node.js subprocess.
func (b *NodeBridge) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cmd != nil && b.cmd.Process != nil {
		return b.cmd.Process.Kill()
	}
	return nil
}

// SendMessage sends a message to the Node.js bridge (stub — real impl would use IPC/gRPC).
func (b *NodeBridge) SendMessage(req *channelv1.SendMessageRequest) error {
	// TODO: implement IPC or gRPC communication with Node.js bridge
	return nil
}

// OnMessage returns the channel of inbound messages from the Node.js bridge.
func (b *NodeBridge) OnMessage() <-chan *channelv1.InboundMessageProto {
	return b.msgCh
}
