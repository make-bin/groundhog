package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	plugin "github.com/make-bin/groundhog/pkg/domain/plugin"
	"github.com/make-bin/groundhog/pkg/domain/plugin/aggregate/plugin_instance"
	"github.com/make-bin/groundhog/pkg/domain/plugin/repository"
	"github.com/make-bin/groundhog/pkg/domain/plugin/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/datastore"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/mapper"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

type pluginRepositoryImpl struct {
	DataStore datastore.DataStore `inject:"datastore"`
}

// NewPluginRepository creates a new PluginRepository implementation.
func NewPluginRepository() repository.PluginRepository {
	return &pluginRepositoryImpl{}
}

// FindAll retrieves all PluginInstance aggregates.
func (r *pluginRepositoryImpl) FindAll(ctx context.Context) ([]*plugin_instance.PluginInstance, error) {
	var pluginPOs []po.PluginPO
	if err := r.DataStore.DB().WithContext(ctx).Find(&pluginPOs).Error; err != nil {
		return nil, err
	}
	plugins := make([]*plugin_instance.PluginInstance, 0, len(pluginPOs))
	for i := range pluginPOs {
		p, err := mapper.PluginPOToDomain(&pluginPOs[i])
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, p)
	}
	return plugins, nil
}

// FindByID retrieves a PluginInstance aggregate by its PluginID.
func (r *pluginRepositoryImpl) FindByID(ctx context.Context, id vo.PluginID) (*plugin_instance.PluginInstance, error) {
	var pluginPO po.PluginPO
	result := r.DataStore.DB().WithContext(ctx).
		Where("plugin_id = ?", id.Value()).
		First(&pluginPO)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, plugin.ErrPluginNotFound
		}
		return nil, result.Error
	}
	return mapper.PluginPOToDomain(&pluginPO)
}

// Save persists a PluginInstance aggregate (upsert by plugin_id).
func (r *pluginRepositoryImpl) Save(ctx context.Context, p *plugin_instance.PluginInstance) error {
	pluginPO, err := mapper.DomainToPluginPO(p)
	if err != nil {
		return err
	}
	return r.DataStore.DB().WithContext(ctx).
		Where("plugin_id = ?", pluginPO.PluginID).
		Save(pluginPO).Error
}

// Delete removes a PluginInstance aggregate by its PluginID.
func (r *pluginRepositoryImpl) Delete(ctx context.Context, id vo.PluginID) error {
	return r.DataStore.DB().WithContext(ctx).
		Where("plugin_id = ?", id.Value()).
		Delete(&po.PluginPO{}).Error
}
