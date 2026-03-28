package plugin

import "errors"

var (
	ErrPluginNotFound    = errors.New("plugin not found")
	ErrPluginNotRunning  = errors.New("plugin is not running")
	ErrPluginMaxRestarts = errors.New("plugin exceeded maximum restart attempts")
	ErrPluginStartFailed = errors.New("plugin failed to start")
)
