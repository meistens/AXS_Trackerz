package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New() *Logger {
	// Create structured logger with JSON output
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true, // Adds file and line number
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{
		logger: logger,
	}
}

func NewWithLevel(level slog.Level) *Logger {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{
		logger: logger,
	}
}

// Info logs info level messages
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Error logs error level messages
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// Warn logs warning level messages
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Debug logs debug level messages
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// With creates a new logger with additional context
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

// WithGroup creates a new logger with a group
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		logger: l.logger.WithGroup(name),
	}
}
