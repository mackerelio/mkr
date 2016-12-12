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
				&mkr.Alert{ID: "2tZhm", Type: "connectivity", Status: "CRITICAL", HostID: "3XYyG", MonitorID: "5rXR3", OpenedAt: 100},
				&mkr.Host{ID: "3XYyG", Name: "app.example.com", Roles: mkr.Roles{"foo": {"bar", "baz"}}, Status: "working"},
				&mkr.MonitorConnectivity{ID: "5rXR3", Type: "connectivity", Name: "connectivity"},
			},
			"2tZhm 1970-01-01 00:01:40 CRITICAL connectivity app.example.com working [foo:bar,baz]",
		},
		{
			&alertSet{
				&mkr.Alert{ID: "2tZhm", Type: "host", Status: "CRITICAL", HostID: "3XYyG", MonitorID: "5rXR3", Value: 15.7, OpenedAt: 200},
				&mkr.Host{ID: "3XYyG", Name: "app.example.com", Roles: mkr.Roles{"foo": {"bar", "baz"}}, Status: "working"},
				&mkr.MonitorHostMetric{ID: "5rXR3", Type: "host", Name: "All::loadavg5", Metric: "loadavg5", Warning: 8.0, Critical: 12.0, Operator: ">"},
			},
			"2tZhm 1970-01-01 00:03:20 CRITICAL All::loadavg5 loadavg5 15.70 > 12.00 app.example.com working [foo:bar,baz]",
		},
		{
			&alertSet{
				&mkr.Alert{ID: "2tZhm", Type: "service", Status: "WARNING", MonitorID: "5rXR3", Value: 15.7, OpenedAt: 300},
				nil,
				&mkr.MonitorServiceMetric{ID: "5rXR3", Type: "service", Service: "ServiceFoo", Name: "bar.baz monitor", Metric: "custom.bar.baz", Warning: 10.0, Critical: 20.0, Operator: ">"},
			},
			"2tZhm 1970-01-01 00:05:00 WARNING bar.baz monitor ServiceFoo custom.bar.baz 15.70 > 10.00",
		},
		{
			&alertSet{
				&mkr.Alert{ID: "2tZhm", Type: "external", Status: "CRITICAL", MonitorID: "5rXR3", Value: 2500, Message: "200", OpenedAt: 400},
				nil,
				&mkr.MonitorExternalHTTP{ID: "5rXR3", Type: "external", Name: "Example Domain", URL: "https://example.com", ResponseTimeWarning: 500, ResponseTimeCritical: 1000, ResponseTimeDuration: 5},
			},
			"2tZhm 1970-01-01 00:06:40 CRITICAL Example Domain https://example.com 2500.00 > 1000.00 msec, status:200",
		},
		{
			&alertSet{
				&mkr.Alert{ID: "2tZhm", Type: "expression", Status: "WARNING", MonitorID: "5rXR3", Value: 15.7, OpenedAt: 500},
				nil,
				&mkr.MonitorExpression{ID: "5rXR3", Type: "expression", Name: "Max loadavg5 monitor", Expression: "max(  \n  roleSlots(  \n    'service:role',\n    'loadavg5'\n  )\n)\n", Warning: 10.0, Critical: 20.0, Operator: ">"},
			},
			"2tZhm 1970-01-01 00:08:20 WARNING Max loadavg5 monitor max(roleSlots('service:role', 'loadavg5')) 15.70 > 10.00",
		},
	}

	for _, testCase := range testCases {
		str := formatJoinedAlert(testCase.alertSet, false)
		if str != testCase.want {
			t.Errorf("should be '%s' but got '%s'", testCase.want, str)
		}
	}
}
