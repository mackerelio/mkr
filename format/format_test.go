package format

import (
	"testing"
	"time"
)

func TestISO8601Extended(t *testing.T) {
	now := time.Now()
	expect := now.Format("2006-01-02T15:04:05-07:00") // ISO 8601 extended format
	got := ISO8601Extended(now)
	if got != expect {
		t.Errorf("should be %q got %q", expect, got)
	}
}
