package mapper

import (
	"encoding/json"
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/plugin/aggregate/plugin_instance"
	"github.com/make-bin/groundhog/pkg/domain/plugin/entity"
	"github.com/make-bin/groundhog/pkg/domain/plugin/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

// DomainToPluginPO converts a PluginInstance aggregate to a PluginPO.
func DomainToPluginPO(p *plugin_instance.PluginInstance) (*po.PluginPO, error) {
	caps := make([]string, 0, len(p.Capabilities()))
	for _, c := range p.Capabilities() {
		caps = append(caps, c.Value())
	}
	capsJSON, err := json.Marshal(caps)
	if err != nil {
		return nil, fmt.Errorf("marshal capabilities: %w", err)
	}
	var name, version, pluginType, entryPoint string
	if m := p.Manifest(); m != nil {
		name = m.Name()
		version = m.Version()
		pluginType = m.PluginType()
		entryPoint = m.EntryPoint()
	}
	return &po.PluginPO{
		PluginID:     p.ID().Value(),
		Name:         name,
		Version:      version,
		PluginType:   pluginType,
		EntryPoint:   entryPoint,
		Status:       int(p.Status()),
		Capabilities: string(capsJSON),
		StartedAt:    p.StartedAt(),
		RestartCount: p.RestartCount(),
	}, nil
}

// PluginPOToDomain converts a PluginPO to a PluginInstance aggregate.
func PluginPOToDomain(p *po.PluginPO) (*plugin_instance.PluginInstance, error) {
	pluginID, err := vo.NewPluginID(p.PluginID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct plugin_id: %w", err)
	}
	manifest := entity.NewPluginManifest(p.Name, p.Version, p.PluginType, p.EntryPoint, nil)
	var caps []vo.Capability
	if p.Capabilities != "" {
		var capStrs []string
		if err := json.Unmarshal([]byte(p.Capabilities), &capStrs); err == nil {
			for _, s := range capStrs {
				if c, err := vo.NewCapability(s); err == nil {
					caps = append(caps, c)
				}
			}
		}
	}
	return plugin_instance.ReconstructPluginInstance(
		pluginID, manifest, nil,
		plugin_instance.PluginStatus(p.Status),
		caps, p.StartedAt, p.RestartCount,
	), nil
}
