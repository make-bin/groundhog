// @AI_GENERATED
package migration

import (
	"fmt"
	"regexp"
	"time"

	"github.com/make-bin/groundhog/pkg/utils/config"
)

// MigrationConfig holds validated migration configuration.
type MigrationConfig struct {
	Enabled     bool
	SourceType  string
	SourcePath  string
	DatabaseURL string
	LockTimeout time.Duration
	TableName   string
}

// BuildMigrationConfig constructs a MigrationConfig from AppConfig.
func BuildMigrationConfig(cfg *config.AppConfig) (*MigrationConfig, error) {
	mc := &MigrationConfig{
		Enabled:     cfg.Migration.Enabled,
		SourceType:  cfg.Migration.SourceType,
		SourcePath:  cfg.Migration.SourcePath,
		LockTimeout: cfg.Migration.LockTimeout,
		TableName:   cfg.Migration.TableName,
	}

	mc.DatabaseURL = BuildDatabaseURL(&cfg.Database)

	if !mc.Enabled {
		return mc, nil
	}

	if err := mc.Validate(); err != nil {
		return nil, err
	}

	return mc, nil
}

// BuildDatabaseURL constructs a PostgreSQL connection URL from DatabaseConfig.
func BuildDatabaseURL(db *config.DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.DBName, db.SSLMode)
}

// Validate checks that the MigrationConfig fields are valid.
func (mc *MigrationConfig) Validate() error {
	if mc.SourceType != "filesystem" && mc.SourceType != "embedded" {
		return fmt.Errorf("invalid migration source_type %q: must be \"filesystem\" or \"embedded\"", mc.SourceType)
	}
	if mc.SourcePath == "" {
		return fmt.Errorf("migration source_path must not be empty")
	}
	if mc.DatabaseURL == "" {
		return fmt.Errorf("migration database URL must not be empty")
	}
	if mc.LockTimeout <= 0 {
		return fmt.Errorf("migration lock_timeout must be positive")
	}
	return nil
}

// passwordPattern matches the password portion of a PostgreSQL URL.
var passwordPattern = regexp.MustCompile(`(postgres://[^:]+:)([^@]+)(@)`)

// SanitizeURL replaces the password in a database URL with ****.
func SanitizeURL(url string) string {
	return passwordPattern.ReplaceAllString(url, "${1}****${3}")
}

// @AI_GENERATED: end
