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

	testCases := []struct {
		alertSet *alertSet
		want     string
	}{
		{
			&alertSet{
				&mkr.Alert{ID: "2tZhm", Type: "connectivity", Status: "critical", HostID: "3XYyG", MonitorID: "5rXR3", OpenedAt: 100},
				&mkr.Host{ID: "3XYyG", Name: "foo", Roles: mkr.Roles{}, Status: "working"},
				&mkr.MonitorConnectivity{ID: "5rXR3", Type: "connectivity"},
			},
			"2tZhm 1970-01-01 00:01:40 critical connectivity foo working []",
		},
	}

	for _, testCase := range testCases {
		str := formatJoinedAlert(testCase.alertSet, false)
		if str != testCase.want {
			t.Errorf("should be '%s' but '%s'", testCase.want, str)
		}
	}
}
