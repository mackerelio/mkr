package apm

import (
	"fmt"
	"io"
)

type testLogger struct {
	w io.Writer
}

func (l *testLogger) Log(prefix, message string) {
	fmt.Fprintln(l.w, prefix, message)
}

func (l *testLogger) Error(err error) {
	fmt.Fprintln(l.w, err.Error())
}
