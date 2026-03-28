package plugin_instance

import (
	"time"

	"github.com/make-bin/groundhog/pkg/domain/plugin/entity"
	"github.com/make-bin/groundhog/pkg/domain/plugin/vo"
)

// PluginStatus represents the lifecycle state of a plugin instance.
type PluginStatus int

const (
	PluginStatusDiscovered PluginStatus = iota
	PluginStatusStarting
	PluginStatusRunning
	PluginStatusStopped
	PluginStatusCrashed
)

// PluginInstance is the aggregate root for a plugin's lifecycle.
type PluginInstance struct {
	id           vo.PluginID
	manifest     *entity.PluginManifest
	process      *entity.PluginProcess
	status       PluginStatus
	capabilities []vo.Capability
	startedAt    *time.Time
	restartCount int
}

// NewPluginInstance creates a new PluginInstance in the Discovered state.
func NewPluginInstance(id vo.PluginID, manifest *entity.PluginManifest) *PluginInstance {
	return &PluginInstance{
		id:           id,
		manifest:     manifest,
		process:      nil,
		status:       PluginStatusDiscovered,
		capabilities: []vo.Capability{},
		startedAt:    nil,
		restartCount: 0,
	}
}

// ReconstructPluginInstance reconstructs a PluginInstance from persisted data.
// This should only be used by repository implementations.
func ReconstructPluginInstance(
	id vo.PluginID,
	manifest *entity.PluginManifest,
	process *entity.PluginProcess,
	status PluginStatus,
	capabilities []vo.Capability,
	startedAt *time.Time,
	restartCount int,
) *PluginInstance {
	return &PluginInstance{
		id:           id,
		manifest:     manifest,
		process:      process,
		status:       status,
		capabilities: capabilities,
		startedAt:    startedAt,
		restartCount: restartCount,
	}
}

// ID returns the plugin instance identifier.
func (p *PluginInstance) ID() vo.PluginID { return p.id }

// Manifest returns the plugin manifest.
func (p *PluginInstance) Manifest() *entity.PluginManifest { return p.manifest }

// Process returns the current plugin process, or nil if not running.
func (p *PluginInstance) Process() *entity.PluginProcess { return p.process }

// Status returns the current plugin status.
func (p *PluginInstance) Status() PluginStatus { return p.status }

// Capabilities returns a copy of the plugin's declared capabilities.
func (p *PluginInstance) Capabilities() []vo.Capability {
	result := make([]vo.Capability, len(p.capabilities))
	copy(result, p.capabilities)
	return result
}

// StartedAt returns the time the plugin was last started, or nil if never started.
func (p *PluginInstance) StartedAt() *time.Time { return p.startedAt }

// RestartCount returns the number of times the plugin has been restarted after a crash.
func (p *PluginInstance) RestartCount() int { return p.restartCount }

// Start transitions the plugin to Running, attaching the given process.
func (p *PluginInstance) Start(process *entity.PluginProcess) {
	p.process = process
	p.status = PluginStatusRunning
	now := time.Now()
	p.startedAt = &now
}

// Stop transitions the plugin to Stopped and clears the process.
func (p *PluginInstance) Stop() {
	p.status = PluginStatusStopped
	p.process = nil
}

// HealthCheck returns true if the plugin is currently Running.
func (p *PluginInstance) HealthCheck() bool {
	return p.status == PluginStatusRunning
}

// RecordCrash increments the restart counter and transitions the plugin to Crashed.
func (p *PluginInstance) RecordCrash() {
	p.restartCount++
	p.status = PluginStatusCrashed
}

// SetCapabilities replaces the plugin's capability list.
func (p *PluginInstance) SetCapabilities(caps []vo.Capability) {
	result := make([]vo.Capability, len(caps))
	copy(result, caps)
	p.capabilities = result
}
