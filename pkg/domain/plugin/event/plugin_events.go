package event

import (
	"time"

	"github.com/make-bin/groundhog/pkg/domain/plugin/vo"
)

// DomainEvent is the interface all plugin domain events must implement.
type DomainEvent interface {
	OccurredAt() time.Time
	EventType() string
}

// PluginStarted is raised when a plugin instance successfully starts.
type PluginStarted struct {
	PluginID   vo.PluginID
	occurredAt time.Time
}

// NewPluginStarted creates a new PluginStarted event.
func NewPluginStarted(pluginID vo.PluginID) PluginStarted {
	return PluginStarted{
		PluginID:   pluginID,
		occurredAt: time.Now(),
	}
}

func (e PluginStarted) OccurredAt() time.Time { return e.occurredAt }
func (e PluginStarted) EventType() string     { return "plugin.instance.started" }

// PluginStopped is raised when a plugin instance stops.
type PluginStopped struct {
	PluginID   vo.PluginID
	occurredAt time.Time
}

// NewPluginStopped creates a new PluginStopped event.
func NewPluginStopped(pluginID vo.PluginID) PluginStopped {
	return PluginStopped{
		PluginID:   pluginID,
		occurredAt: time.Now(),
	}
}

func (e PluginStopped) OccurredAt() time.Time { return e.occurredAt }
func (e PluginStopped) EventType() string     { return "plugin.instance.stopped" }

// PluginCrashed is raised when a plugin instance crashes unexpectedly.
type PluginCrashed struct {
	PluginID     vo.PluginID
	RestartCount int
	occurredAt   time.Time
}

// NewPluginCrashed creates a new PluginCrashed event.
func NewPluginCrashed(pluginID vo.PluginID, restartCount int) PluginCrashed {
	return PluginCrashed{
		PluginID:     pluginID,
		RestartCount: restartCount,
		occurredAt:   time.Now(),
	}
}

func (e PluginCrashed) OccurredAt() time.Time { return e.occurredAt }
func (e PluginCrashed) EventType() string     { return "plugin.instance.crashed" }
