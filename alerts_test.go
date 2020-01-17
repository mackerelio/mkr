package main

import (
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"
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
				&mackerel.Alert{ID: "2tZhm", Type: "connectivity", Status: "CRITICAL", HostID: "3XYyG", MonitorID: "5rXR3", OpenedAt: 100},
				&mackerel.Host{ID: "3XYyG", Name: "app.example.com", Roles: mackerel.Roles{"foo": {"bar", "baz"}}, Status: "working"},
				&mackerel.MonitorConnectivity{ID: "5rXR3", Type: "connectivity", Name: "connectivity"},
			},
			"2tZhm 1970-01-01 00:01:40 CRITICAL connectivity app.example.com working [foo:bar,baz]",
		},
		{
			&alertSet{
				&mackerel.Alert{ID: "2tZhm", Type: "host", Status: "CRITICAL", HostID: "3XYyG", MonitorID: "5rXR3", Value: 15.7, OpenedAt: 200},
				&mackerel.Host{ID: "3XYyG", Name: "app.example.com", Roles: mackerel.Roles{"foo": {"bar", "baz"}}, Status: "working"},
				&mackerel.MonitorHostMetric{ID: "5rXR3", Type: "host", Name: "All::loadavg5", Metric: "loadavg5", Warning: pfloat64(8.0), Critical: pfloat64(12.0), Operator: ">"},
			},
			"2tZhm 1970-01-01 00:03:20 CRITICAL All::loadavg5 loadavg5 15.70 > 12.00 app.example.com working [foo:bar,baz]",
		},
		{
			&alertSet{
				&mackerel.Alert{ID: "2tZhm", Type: "service", Status: "WARNING", MonitorID: "5rXR3", Value: 15.7, OpenedAt: 300},
				nil,
				&mackerel.MonitorServiceMetric{ID: "5rXR3", Type: "service", Service: "ServiceFoo", Name: "bar.baz monitor", Metric: "custom.bar.baz", Warning: pfloat64(10.0), Critical: pfloat64(20.0), Operator: ">"},
			},
			"2tZhm 1970-01-01 00:05:00 WARNING bar.baz monitor ServiceFoo custom.bar.baz 15.70 > 10.00",
		},
		{
			&alertSet{
				&mackerel.Alert{ID: "2tZhm", Type: "external", Status: "CRITICAL", MonitorID: "5rXR3", Value: 2500, Message: "200", OpenedAt: 400},
				nil,
				&mackerel.MonitorExternalHTTP{ID: "5rXR3", Type: "external", Name: "Example Domain", URL: "https://example.com", ResponseTimeWarning: pfloat64(500), ResponseTimeCritical: pfloat64(1000), ResponseTimeDuration: puint64(5)},
			},
			"2tZhm 1970-01-01 00:06:40 CRITICAL Example Domain https://example.com 2500.00 > 1000.00 msec, status:200",
		},
		{
			&alertSet{
				&mackerel.Alert{ID: "2tZhm", Type: "expression", Status: "WARNING", MonitorID: "5rXR3", Value: 15.7, OpenedAt: 500},
				nil,
				&mackerel.MonitorExpression{ID: "5rXR3", Type: "expression", Name: "Max loadavg5 monitor", Expression: "max(  \n  roleSlots(  \n    'service:role',\n    'loadavg5'\n  )\n)\n", Warning: pfloat64(10.0), Critical: pfloat64(20.0), Operator: ">"},
			},
			"2tZhm 1970-01-01 00:08:20 WARNING Max loadavg5 monitor max(roleSlots('service:role', 'loadavg5')) 15.70 > 10.00",
		},
		{
			&alertSet{
				&mackerel.Alert{ID: "2tZhm", Type: "check", Status: "WARNING", MonitorID: "5rXR3", OpenedAt: 500, Message: "Short check monitoring description"},
				&mackerel.Host{ID: "3XYyG", Name: "app.example.com", Roles: mackerel.Roles{"foo": {"bar", "baz"}}, Status: "working"},
				nil,
			},
			"2tZhm 1970-01-01 00:08:20 WARNING Short check monitoring description app.example.com working [foo:bar,baz]",
		},
		{
			&alertSet{
				&mackerel.Alert{ID: "2tZhm", Type: "check", Status: "WARNING", MonitorID: "5rXR3", OpenedAt: 500, Message: "description with\nnew line\nmore new line "},
				&mackerel.Host{ID: "3XYyG", Name: "app.example.com", Roles: mackerel.Roles{"foo": {"bar", "baz"}}, Status: "working"},
				nil,
			},
			"2tZhm 1970-01-01 00:08:20 WARNING description with... app.example.com working [foo:bar,baz]",
		},
		{
			&alertSet{
				&mackerel.Alert{ID: "2tZhm", Type: "check", Status: "WARNING", MonitorID: "5rXR3", OpenedAt: 500, Message: "long long long long long long long long long long long long long long long long long long long そして長い長い alert"},
				&mackerel.Host{ID: "3XYyG", Name: "app.example.com", Roles: mackerel.Roles{"foo": {"bar", "baz"}}, Status: "working"},
				nil,
			},
			"2tZhm 1970-01-01 00:08:20 WARNING long long long long long long long long long long long long long long long long long long long そして長い... app.example.com working [foo:bar,baz]",
		},
		{
			&alertSet{
				&mackerel.Alert{ID: "2tZhm", Type: "anomalyDetection", Status: "WARNING", MonitorID: "5rXR3", OpenedAt: 500},
				&mackerel.Host{ID: "3XYyG", Name: "app.example.com", Roles: mackerel.Roles{"foo": {"bar", "baz"}}, Status: "working"},
				&mackerel.MonitorAnomalyDetection{ID: "5rXR3", Type: "anomalyDetection", Name: "My anomaly detection for roles", WarningSensitivity: "insensitive", MaxCheckAttempts: 5, Scopes: []string{"foo: bar", "foo: baz"}},
			},
			"2tZhm 1970-01-01 00:08:20 WARNING My anomaly detection for roles app.example.com working [foo:bar,baz]",
		},
	}

	for _, testCase := range testCases {
		str := formatJoinedAlert(testCase.alertSet, false)
		if str != testCase.want {
			t.Errorf("should be '%s' but got '%s'", testCase.want, str)
		}
	}
}
