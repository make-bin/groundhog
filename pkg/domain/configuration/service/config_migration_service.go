// @AI_GENERATED
package service

import (
	"github.com/make-bin/groundhog/pkg/domain/configuration/aggregate/config_root"
)

// ConfigMigrationService migrates a ConfigRoot from an older schema version to the current version.
type ConfigMigrationService interface {
	Migrate(root *config_root.ConfigRoot) (*config_root.ConfigRoot, error)
}

// @AI_GENERATED: end
