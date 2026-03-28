// @AI_GENERATED
package logger

// Logger defines the structured logging interface used across all layers.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, console
}

// @AI_GENERATED: end
