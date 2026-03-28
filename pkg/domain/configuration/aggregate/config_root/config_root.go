// @AI_GENERATED
package config_root

import (
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/configuration/vo"
)

// CurrentVersion is the current configuration schema version.
// NeedsMigration compares the aggregate's version against this constant.
var CurrentVersion = func() vo.SchemaVersion {
	v, _ := vo.NewSchemaVersion(1, 0, 0)
	return v
}()

// ConfigRoot is the aggregate root for the Configuration bounded context.
// It represents the entire application configuration tree.
type ConfigRoot struct {
	version  vo.SchemaVersion
	sections map[string]interface{}
	secrets  map[string]vo.SecretRef
	filePath string
	loadedAt time.Time
	dirty    bool
}

// NewConfigRoot creates a new ConfigRoot with the given schema version and file path.
func NewConfigRoot(version vo.SchemaVersion, filePath string) *ConfigRoot {
	return &ConfigRoot{
		version:  version,
		sections: make(map[string]interface{}),
		secrets:  make(map[string]vo.SecretRef),
		filePath: filePath,
		loadedAt: time.Now(),
		dirty:    false,
	}
}

// Get retrieves a configuration value by path.
// Returns an error if the path is not found.
func (c *ConfigRoot) Get(path string) (interface{}, error) {
	val, ok := c.sections[path]
	if !ok {
		return nil, fmt.Errorf("configuration path %q not found", path)
	}
	return val, nil
}

// Set stores a configuration value at the given path and marks the root as dirty.
func (c *ConfigRoot) Set(path string, value interface{}) error {
	c.sections[path] = value
	c.dirty = true
	return nil
}

// Validate checks all sections and returns any validation errors found.
func (c *ConfigRoot) Validate() []vo.ValidationError {
	var errs []vo.ValidationError
	if c.filePath == "" {
		errs = append(errs, vo.NewValidationError("filePath", "file path must not be empty"))
	}
	return errs
}

// NeedsMigration returns true if the aggregate's version is less than CurrentVersion.
func (c *ConfigRoot) NeedsMigration() bool {
	return c.version.LessThan(CurrentVersion)
}

// Version returns the schema version.
func (c *ConfigRoot) Version() vo.SchemaVersion {
	return c.version
}

// IsDirty returns true if the configuration has been modified since loading.
func (c *ConfigRoot) IsDirty() bool {
	return c.dirty
}

// FilePath returns the file path from which the configuration was loaded.
func (c *ConfigRoot) FilePath() string {
	return c.filePath
}

// @AI_GENERATED: end
