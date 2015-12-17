package main

import (
	"testing"

	mkr "github.com/mackerelio/mackerel-client-go"
)

func TestFormatJoinedAlert(t *testing.T) {
	a := mkr.Alert{}
	h := mkr.Host{}
	m := mkr.Monitor{}
	as := alertSet{&a, &h, &m}
	answer := " 1970-01-01 09:00:00     []"

	str := formatJoinedAlert(&as, false)
	if str != answer {
		t.Errorf("should be '%s' but '%s'", answer, str)
	}
}
