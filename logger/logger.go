package logger

// We borrow this code from github.com/motemen/ghq/utils

import (
	"fmt"
	"os"

	colorine "github.com/motemen/go-colorine"
)

var logger = &colorine.Logger{
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

func init() {
	logger.SetOutput(os.Stderr)
}

// Log outputs `message` with `prefix` by go-colorine
func Log(prefix, message string) {
	logger.Log(prefix, message)
}

// Logf outputs `message` with `prefix` by go-colorine
func Logf(prefix, message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	logger.Log(prefix, msg)
}

// ErrorIf outputs log if `err` occurs.
func ErrorIf(err error) bool {
	if err != nil {
		Log("error", err.Error())
		return true
	}

	return false
}

// DieIf outputs log and exit(1) if `err` occurs.
func DieIf(err error) {
	if err != nil {
		Log("error", err.Error())
		os.Exit(1)
	}
}

// PanicIf raise panic if `err` occurs.
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
