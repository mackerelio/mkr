package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestQuietMode(t *testing.T) {
	var buf bytes.Buffer
	l := New()
	l.logger.SetOutput(&buf)

	// Info-level prefixes should be logged normally
	l.Log("thrown", "metric1")
	if !strings.Contains(buf.String(), "metric1") {
		t.Error("expected 'metric1' in output")
	}

	buf.Reset()

	// Enable quiet mode: info-level should be suppressed
	l.SetQuiet(true)
	for _, prefix := range []string{"", "info", "created", "updated", "thrown"} {
		l.Log(prefix, "should-not-appear")
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output in quiet mode for info-level prefixes, got: %s", buf.String())
	}

	// Warning and error should still appear
	l.Log("warning", "warn-msg")
	if !strings.Contains(buf.String(), "warn-msg") {
		t.Error("expected 'warn-msg' in output even in quiet mode")
	}

	buf.Reset()
	l.Log("error", "err-msg")
	if !strings.Contains(buf.String(), "err-msg") {
		t.Error("expected 'err-msg' in output even in quiet mode")
	}
}
