package logger

// Originally from github.com/motemen/ghq/utils

import (
	"fmt"
	"os"

	colorine "github.com/motemen/go-colorine"
)

// Logger is wrapped go-colorine logger for mkr
type Logger struct {
	logger *colorine.Logger
}

// New is constructor for new colorine logger
func New() *Logger {
	logger := &colorine.Logger{
		Prefixes: colorine.Prefixes{
			"warning": colorine.Warn,

			"error": colorine.Error,

			"":        colorine.Info,
			"info":    colorine.Info,
			"created": colorine.Info,
			"updated": colorine.Info,
			"thrown":  colorine.Info,
			"retired": colorine.Info,
		},
	}

	// Default output
	logger.SetOutput(os.Stderr)
	return &Logger{logger: logger}
}

// Log outputs `message` with `prefix` by go-colorine
func (l *Logger) Log(prefix, message string) {
	l.logger.Log(prefix, message)
}

// Logf outputs `message` with `prefix` by go-colorine
func (l *Logger) Logf(prefix, message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	l.logger.Log(prefix, msg)
}

// Error outputs log given non-nil `err`
func (l *Logger) Error(err error) {
	l.Log("error", err.Error())
}

var defaultLogger = New()

// Log outputs `message` with `prefix` by go-colorine
func Log(prefix, message string) {
	defaultLogger.Log(prefix, message)
}

// Logf outputs `message` with `prefix` by go-colorine
func Logf(prefix, message string, args ...interface{}) {
	defaultLogger.Logf(prefix, message, args...)
}

// ErrorIf outputs log if `err` occurs.
func ErrorIf(err error) bool {
	if err == nil {
		return false
	}
	defaultLogger.Error(err)
	return true
}

// DieIf outputs log and exit(1) if `err` occurs.
func DieIf(err error) {
	if ErrorIf(err) {
		os.Exit(1)
	}
}
