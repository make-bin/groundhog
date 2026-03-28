package service

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/plugin/aggregate/plugin_instance"
	"github.com/make-bin/groundhog/pkg/domain/plugin/vo"
)

// PluginLifecycleService manages the lifecycle of plugin instances.
type PluginLifecycleService interface {
	Discover(ctx context.Context) ([]*plugin_instance.PluginInstance, error)
	Start(ctx context.Context, id vo.PluginID) error
	Stop(ctx context.Context, id vo.PluginID) error
	Restart(ctx context.Context, id vo.PluginID) error
	HealthCheck(ctx context.Context, id vo.PluginID) (bool, error)
}
