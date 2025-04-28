// logger/logger.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the basic logging interface.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// zapLogger is an implementation of Logger using Zap's SugaredLogger.
type zapLogger struct {
	logger *zap.SugaredLogger
}

// NewLogger creates and configures a Zap SugaredLogger with:
//   - ISO8601 timestamp format
//   - colored level output
//   - short caller information (file:line)
//
// It returns the logger wrapped in our Logger interface.
func NewLogger() (Logger, error) {
	cfg := zap.NewDevelopmentConfig()

	// Use ISO8601 for timestamps
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// Colorize level (INFO, ERROR, etc.) for readability
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// Show caller as short file:line
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Build the logger, adding caller and skipping one frame (this function)
	zapLog, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: zapLog.Sugar()}, nil
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
