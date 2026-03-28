package entity

// PluginManifest describes a plugin's metadata and capabilities.
type PluginManifest struct {
	name         string
	version      string
	pluginType   string
	entryPoint   string
	dependencies []string
}

// NewPluginManifest creates a new PluginManifest with the given metadata.
func NewPluginManifest(name, version, pluginType, entryPoint string, dependencies []string) *PluginManifest {
	deps := make([]string, len(dependencies))
	copy(deps, dependencies)
	return &PluginManifest{
		name:         name,
		version:      version,
		pluginType:   pluginType,
		entryPoint:   entryPoint,
		dependencies: deps,
	}
}

// Name returns the plugin name.
func (m *PluginManifest) Name() string { return m.name }

// Version returns the plugin version.
func (m *PluginManifest) Version() string { return m.version }

// PluginType returns the plugin type identifier.
func (m *PluginManifest) PluginType() string { return m.pluginType }

// EntryPoint returns the plugin entry point path.
func (m *PluginManifest) EntryPoint() string { return m.entryPoint }

// Dependencies returns a copy of the plugin dependency list.
func (m *PluginManifest) Dependencies() []string {
	result := make([]string, len(m.dependencies))
	copy(result, m.dependencies)
	return result
}
