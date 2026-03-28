// @AI_GENERATED
package migration

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// Register file source driver for golang-migrate.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/make-bin/groundhog/pkg/utils/config"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

const (
	maxRetries   = 3
	initialDelay = 2 * time.Second
)

// RunMigrations executes database migrations on application startup.
func RunMigrations(cfg *config.AppConfig, log logger.Logger) error {
	if cfg == nil {
		return fmt.Errorf("migration: config must not be nil")
	}
	if log == nil {
		return fmt.Errorf("migration: logger must not be nil")
	}

	mc, err := BuildMigrationConfig(cfg)
	if err != nil {
		return fmt.Errorf("migration: invalid config: %w", err)
	}

	if !mc.Enabled {
		log.Info("database migrations are disabled, skipping")
		return nil
	}

	sanitizedURL := SanitizeURL(mc.DatabaseURL)
	log.Info("starting database migrations",
		"source_type", mc.SourceType,
		"source_path", mc.SourcePath,
		"database", sanitizedURL,
	)

	// Resolve source URL.
	sourceURL, err := resolveSourceURL(mc)
	if err != nil {
		return fmt.Errorf("migration: %w", err)
	}

	// Create migration manager with retry on connection errors.
	mgr, err := createManagerWithRetry(sourceURL, mc.DatabaseURL, log)
	if err != nil {
		return fmt.Errorf("migration: failed to connect after %d attempts: %w", maxRetries, err)
	}
	defer func() {
		if closeErr := mgr.Close(); closeErr != nil {
			log.Warn("failed to close migration manager", "error", closeErr)
		}
	}()

	// Check for dirty state.
	version, dirty, verErr := mgr.Version()
	if verErr != nil && !errors.Is(verErr, migrate.ErrNilVersion) {
		return fmt.Errorf("migration: failed to get version: %w", verErr)
	}
	if dirty {
		return fmt.Errorf(
			"migration: database is in a dirty state at version %d. "+
				"Recovery steps: "+
				"1) If the migration partially succeeded, run: migrate force %d "+
				"2) To rollback, run: migrate down 1 "+
				"3) Check logs and fix the migration manually",
			version, version,
		)
	}

	// Execute Up().
	ctx := context.Background()
	if err := mgr.Up(ctx); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("database migrations already up-to-date", "version", version)
			return nil
		}
		return fmt.Errorf("migration: failed to apply: %w", err)
	}

	// Log final version.
	newVersion, _, _ := mgr.Version()
	log.Info("database migrations completed successfully", "version", newVersion)
	return nil
}

// resolveSourceURL builds the source URL from the MigrationConfig.
func resolveSourceURL(mc *MigrationConfig) (string, error) {
	switch mc.SourceType {
	case "filesystem":
		src := NewFilesystemSource(mc.SourcePath)
		if err := src.Validate(); err != nil {
			return "", err
		}
		return src.GetSourceURL(), nil
	case "embedded":
		return "", fmt.Errorf("embedded source requires explicit setup via NewEmbeddedSource")
	default:
		return "", fmt.Errorf("unsupported source type: %s", mc.SourceType)
	}
}

// createManagerWithRetry attempts to create a MigrationManager with exponential backoff.
func createManagerWithRetry(sourceURL, databaseURL string, log logger.Logger) (MigrationManager, error) {
	var lastErr error
	delay := initialDelay

	for attempt := 1; attempt <= maxRetries; attempt++ {
		mgr, err := NewMigrationManager(sourceURL, databaseURL)
		if err == nil {
			return mgr, nil
		}

		lastErr = err
		if !isConnectionError(err) {
			return nil, err
		}

		sanitizedErr := SanitizeURL(err.Error())
		log.Warn("migration connection failed, retrying",
			"attempt", attempt,
			"max_attempts", maxRetries,
			"delay", delay,
			"error", sanitizedErr,
		)

		time.Sleep(delay)
		delay *= 2
	}

	return nil, lastErr
}

// isConnectionError checks if the error is a transient connection error worth retrying.
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	connectionPatterns := []string{
		"connection refused",
		"dial tcp",
		"no such host",
		"connection reset",
		"connection timed out",
		"i/o timeout",
	}
	for _, pattern := range connectionPatterns {
		if strings.Contains(msg, pattern) {
			return true
		}
	}
	return false
}

// @AI_GENERATED: end
