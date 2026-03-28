package repository

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/plugin/aggregate/plugin_instance"
	"github.com/make-bin/groundhog/pkg/domain/plugin/vo"
)

// PluginRepository defines the data access contract for the PluginInstance aggregate.
type PluginRepository interface {
	FindAll(ctx context.Context) ([]*plugin_instance.PluginInstance, error)
	FindByID(ctx context.Context, id vo.PluginID) (*plugin_instance.PluginInstance, error)
	Save(ctx context.Context, plugin *plugin_instance.PluginInstance) error
	Delete(ctx context.Context, id vo.PluginID) error
}
