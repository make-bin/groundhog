package entity

import "time"

// PluginProcess holds runtime information about a running plugin process.
type PluginProcess struct {
	pid       int
	startedAt time.Time
	grpcAddr  string
}

// NewPluginProcess creates a new PluginProcess with the given runtime details.
func NewPluginProcess(pid int, startedAt time.Time, grpcAddr string) *PluginProcess {
	return &PluginProcess{
		pid:       pid,
		startedAt: startedAt,
		grpcAddr:  grpcAddr,
	}
}

// PID returns the operating system process ID.
func (p *PluginProcess) PID() int { return p.pid }

// StartedAt returns the time the process was started.
func (p *PluginProcess) StartedAt() time.Time { return p.startedAt }

// GRPCAddr returns the gRPC address the plugin is listening on.
func (p *PluginProcess) GRPCAddr() string { return p.grpcAddr }
