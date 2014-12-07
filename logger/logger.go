package logger

// We borrow this code from github.com/motemen/ghq/utils

import (
	"os"

	colorine "github.com/motemen/go-colorine"
)

var logger = &colorine.Logger{
	Prefixes: colorine.Prefixes{
		"warning": colorine.Warn,

		"error": colorine.Error,

		"": colorine.Info,
		"created":  colorine.Info,
		"updated":  colorine.Info,
		"thrown":  colorine.Info,
		"retired":  colorine.Info,
	},
}

func Log(prefix, message string) {
	logger.Log(prefix, message)
}

func ErrorIf(err error) bool {
	if err != nil {
		Log("error", err.Error())
		return true
	}

	return false
}

func DieIf(err error) {
	if err != nil {
		Log("error", err.Error())
		os.Exit(1)
	}
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
