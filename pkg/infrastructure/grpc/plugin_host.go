package grpc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	channelv1 "github.com/make-bin/groundhog/proto/gen/channel/v1"
)

const maxRestarts = 3

// managedPlugin holds the runtime state of a managed plugin subprocess.
type managedPlugin struct {
	id           string
	binaryPath   string
	cmd          *exec.Cmd
	restartCount int
	mu           sync.Mutex
}

// PluginHost manages channel plugin subprocesses.
type PluginHost struct {
	pluginDir string
	plugins   map[string]*managedPlugin
	mu        sync.RWMutex
	msgCh     chan *channelv1.InboundMessageProto
}

// NewPluginHost creates a new PluginHost that discovers plugins in pluginDir.
func NewPluginHost(pluginDir string) *PluginHost {
	return &PluginHost{
		pluginDir: pluginDir,
		plugins:   make(map[string]*managedPlugin),
		msgCh:     make(chan *channelv1.InboundMessageProto, 100),
	}
}

// Discover scans pluginDir for executable binaries and registers them.
func (h *PluginHost) Discover() error {
	entries, err := os.ReadDir(h.pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no plugins dir is OK
		}
		return fmt.Errorf("plugin_host: discover: %w", err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		// Check if executable
		if info.Mode()&0111 == 0 {
			continue
		}
		id := entry.Name()
		binaryPath := filepath.Join(h.pluginDir, entry.Name())
		h.plugins[id] = &managedPlugin{
			id:         id,
			binaryPath: binaryPath,
		}
	}
	return nil
}

// StartAll starts all discovered plugins.
func (h *PluginHost) StartAll(ctx context.Context) error {
	h.mu.RLock()
	ids := make([]string, 0, len(h.plugins))
	for id := range h.plugins {
		ids = append(ids, id)
	}
	h.mu.RUnlock()

	for _, id := range ids {
		if err := h.startPlugin(ctx, id); err != nil {
			return fmt.Errorf("plugin_host: start %s: %w", id, err)
		}
	}
	return nil
}

// startPlugin starts a single plugin subprocess and monitors it.
func (h *PluginHost) startPlugin(ctx context.Context, id string) error {
	h.mu.RLock()
	mp, ok := h.plugins[id]
	h.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %s not found", id)
	}

	mp.mu.Lock()
	defer mp.mu.Unlock()

	cmd := exec.CommandContext(ctx, mp.binaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start subprocess: %w", err)
	}
	mp.cmd = cmd

	// Monitor for crash and auto-restart
	go h.monitorPlugin(ctx, id)
	return nil
}

// monitorPlugin watches a plugin process and restarts it on crash (up to maxRestarts).
func (h *PluginHost) monitorPlugin(ctx context.Context, id string) {
	for {
		h.mu.RLock()
		mp, ok := h.plugins[id]
		h.mu.RUnlock()
		if !ok {
			return
		}

		mp.mu.Lock()
		cmd := mp.cmd
		mp.mu.Unlock()

		if cmd == nil {
			return
		}

		// Wait for process to exit
		_ = cmd.Wait()

		// Check context cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		mp.mu.Lock()
		if mp.restartCount >= maxRestarts {
			mp.mu.Unlock()
			return
		}
		mp.restartCount++
		mp.mu.Unlock()

		// Brief backoff before restart
		time.Sleep(time.Second)

		_ = h.startPlugin(ctx, id)
	}
}

// Stop stops a plugin subprocess by ID.
func (h *PluginHost) Stop(pluginID string) error {
	h.mu.RLock()
	mp, ok := h.plugins[pluginID]
	h.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %s not found", pluginID)
	}

	mp.mu.Lock()
	defer mp.mu.Unlock()
	if mp.cmd != nil && mp.cmd.Process != nil {
		return mp.cmd.Process.Kill()
	}
	return nil
}

// SendMessage sends a message to a plugin. In this stub implementation it is a no-op
// since real gRPC transport is not wired yet.
func (h *PluginHost) SendMessage(pluginID string, req *channelv1.SendMessageRequest) error {
	h.mu.RLock()
	_, ok := h.plugins[pluginID]
	h.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %s not found", pluginID)
	}
	// TODO: wire real gRPC transport
	return nil
}

// OnMessage returns a channel that receives inbound messages from all plugins.
func (h *PluginHost) OnMessage() <-chan *channelv1.InboundMessageProto {
	return h.msgCh
}
