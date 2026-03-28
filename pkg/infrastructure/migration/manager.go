// @AI_GENERATED
package migration

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	// Register postgres database driver for golang-migrate.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

// MigrationManager defines the interface for database migration operations.
type MigrationManager interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
	Version() (version uint, dirty bool, err error)
	Close() error
}

// migrationManager implements MigrationManager using golang-migrate.
type migrationManager struct {
	m *migrate.Migrate
}

// NewMigrationManager creates a new MigrationManager from a source URL and database URL.
func NewMigrationManager(sourceURL, databaseURL string) (MigrationManager, error) {
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}
	return &migrationManager{m: m}, nil
}

// Up applies all pending migrations.
func (mm *migrationManager) Up(_ context.Context) error {
	return mm.m.Up()
}

// Down rolls back the last migration.
func (mm *migrationManager) Down(_ context.Context) error {
	return mm.m.Down()
}

// Version returns the current migration version and dirty state.
func (mm *migrationManager) Version() (uint, bool, error) {
	return mm.m.Version()
}

// Close releases resources held by the migration manager.
func (mm *migrationManager) Close() error {
	srcErr, dbErr := mm.m.Close()
	if srcErr != nil {
		return fmt.Errorf("failed to close migration source: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close migration database: %w", dbErr)
	}
	return nil
}

// @AI_GENERATED: end
