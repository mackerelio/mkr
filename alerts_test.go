package main

import (
	"testing"
	"time"

	mkr "github.com/mackerelio/mackerel-client-go"
)

func TestFormatJoinedAlert(t *testing.T) {

	const location = "UTC"

	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 0*60*60)
	}
	time.Local = loc

	a := &mkr.Alert{ID: "123", Type: "connectivity", Status: "critical", HostID: "1234", MonitorID: "12345", OpenedAt: 100}
	h := &mkr.Host{ID: "1234", Name: "foo", Roles: mkr.Roles{}, Status: "working"}
	m := &mkr.MonitorConnectivity{ID: "12345", Type: "connectivity"}
	as := alertSet{a, h, m}
	answer := "123 1970-01-01 00:01:40 critical connectivity foo working []"

	str := formatJoinedAlert(&as, false)
	if str != answer {
		t.Errorf("should be '%s' but '%s'", answer, str)
	}
}
