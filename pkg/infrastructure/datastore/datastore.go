// @AI_GENERATED
package datastore

import (
	"fmt"

	"github.com/make-bin/groundhog/pkg/utils/config"
	"github.com/make-bin/groundhog/pkg/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DataStore provides a unified database abstraction for repository implementations.
type DataStore interface {
	DB() *gorm.DB
	Close() error
}

type dataStore struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewDataStore initializes a GORM v2 PostgreSQL connection using the provided
// DatabaseConfig and returns a DataStore. It configures the connection pool
// according to the config values and returns a descriptive error on failure.
func NewDataStore(cfg *config.DatabaseConfig, logger logger.Logger) (DataStore, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database at %s:%d/%s: %w", cfg.Host, cfg.Port, cfg.DBName, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	logger.Info("database connection established",
		"host", cfg.Host,
		"port", cfg.Port,
		"dbname", cfg.DBName,
		"max_idle_conns", cfg.MaxIdleConns,
		"max_open_conns", cfg.MaxOpenConns,
		"conn_max_lifetime", cfg.ConnMaxLifetime,
	)

	return &dataStore{db: db, logger: logger}, nil
}

// DB returns the underlying *gorm.DB instance.
func (ds *dataStore) DB() *gorm.DB {
	return ds.db
}

// Close gets the underlying *sql.DB and closes the database connection.
func (ds *dataStore) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB for close: %w", err)
	}
	ds.logger.Info("closing database connection")
	return sqlDB.Close()
}

// @AI_GENERATED: end
