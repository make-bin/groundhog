// @AI_GENERATED
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implements Logger using zap's SugaredLogger.
type zapLogger struct {
	sugar *zap.SugaredLogger
}

// NewLogger creates a new Logger backed by Zap.
func NewLogger(cfg LogConfig) (Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	var zapCfg zap.Config
	switch cfg.Format {
	case "json":
		zapCfg = zap.NewProductionConfig()
	case "console", "":
		zapCfg = zap.NewDevelopmentConfig()
	default:
		return nil, fmt.Errorf("unsupported log format: %s", cfg.Format)
	}

	zapCfg.Level = zap.NewAtomicLevelAt(level)

	z, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return &zapLogger{sugar: z.Sugar()}, nil
}

func parseLevel(lvl string) (zapcore.Level, error) {
	switch lvl {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info", "":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unsupported log level: %s", lvl)
	}
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) With(keysAndValues ...interface{}) Logger {
	return &zapLogger{sugar: l.sugar.With(keysAndValues...)}
}

// @AI_GENERATED: end
