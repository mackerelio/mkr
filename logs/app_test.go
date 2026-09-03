package logs

import (
	"bytes"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/stretchr/testify/assert"
)

func TestLogsApp_findLogs(t *testing.T) {
	traceID := "0123456789abcdef0123456789abcdef"

	testCases := []struct {
		id       string
		resp     *mackerel.FindLogsResponse
		jqFilter string
		expected string
	}{
		{
			id: "default",
			resp: &mackerel.FindLogsResponse{
				Results: []*mackerel.SimpleLog{
					{
						Cursor:             "cursor1",
						Timestamp:          time.Date(2026, 9, 1, 0, 12, 34, 0, time.UTC),
						EffectiveTimestamp: time.Date(2026, 9, 1, 0, 12, 34, 0, time.UTC),
						Severity:           mackerel.LogSeverityError,
						SeverityText:       "Error",
						SeverityNumber:     17,
						Body:               "connection timeout",
						TraceID:            &traceID,
						ServiceName:        "shoppingcart",
						ServiceNamespace:   "shop",
						Attributes:         []*mackerel.LogAttribute{},
						ResourceAttributes: []*mackerel.LogAttribute{},
						ScopeAttributes:    []*mackerel.LogAttribute{},
					},
				},
				PageInfo: mackerel.LogPageInfo{HasNextPage: false, HasPreviousPage: false},
			},
			expected: `{
    "results": [
        {
            "cursor": "cursor1",
            "timestamp": "2026-09-01T00:12:34Z",
            "effectiveTimestamp": "2026-09-01T00:12:34Z",
            "severity": "ERROR",
            "severityText": "Error",
            "severityNumber": 17,
            "body": "connection timeout",
            "traceId": "0123456789abcdef0123456789abcdef",
            "serviceName": "shoppingcart",
            "serviceNamespace": "shop",
            "attributes": [],
            "resourceAttributes": [],
            "scopeAttributes": []
        }
    ],
    "pageInfo": {
        "hasNextPage": false,
        "hasPreviousPage": false
    }
}
`,
		},
		{
			id: "jq",
			resp: &mackerel.FindLogsResponse{
				Results: []*mackerel.SimpleLog{{Body: "hello"}},
			},
			jqFilter: ".results[].body",
			expected: `hello
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			client := mackerelclient.NewMockClient(
				mackerelclient.MockFindLogs(func(param *mackerel.FindLogsParam) (*mackerel.FindLogsResponse, error) {
					return tc.resp, nil
				}),
			)
			out := new(bytes.Buffer)
			app := &logsApp{
				client:    client,
				outStream: out,
				jqFilter:  tc.jqFilter,
			}
			param := &mackerel.FindLogsParam{ServiceName: "shoppingcart"}
			assert.NoError(t, app.findLogs(t.Context(), param))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
