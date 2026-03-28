// @AI_GENERATED
package migration

import (
	"fmt"

	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// migrateLogger adapts the project Logger to golang-migrate's logger.Logger interface.
type migrateLogger struct {
	log logger.Logger
}

// newMigrateLogger creates a new golang-migrate logger adapter.
func newMigrateLogger(log logger.Logger) *migrateLogger {
	return &migrateLogger{log: log}
}

// Printf logs a message at INFO level.
func (l *migrateLogger) Printf(format string, v ...interface{}) {
	l.log.Info(fmt.Sprintf(format, v...))
}

// Verbose returns true to enable verbose logging from golang-migrate.
func (l *migrateLogger) Verbose() bool {
	return true
}

// @AI_GENERATED: end
