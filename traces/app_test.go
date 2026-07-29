package traces

import (
	"bytes"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/stretchr/testify/assert"
)

func TestTracesApp_getTrace(t *testing.T) {
	traceID := "0123456789abcdef0123456789abcdef"

	testCases := []struct {
		id       string
		trace    *mackerel.TraceResponse
		jqFilter string
		expected string
	}{
		{
			id: "default",
			trace: &mackerel.TraceResponse{
				Spans: []*mackerel.Span{
					{
						TraceID:   traceID,
						SpanID:    "0123456780123456",
						Name:      "test-span",
						Kind:      mackerel.SpanKindInternal,
						StartTime: time.Date(2025, 7, 9, 14, 3, 2, 0, time.UTC),
						EndTime:   time.Date(2025, 7, 9, 14, 3, 3, 0, time.UTC),
					},
				},
			},
			expected: `{
    "spans": [
        {
            "traceId": "0123456789abcdef0123456789abcdef",
            "spanId": "0123456780123456",
            "traceState": "",
            "name": "test-span",
            "kind": "internal",
            "startTime": "2025-07-09T14:03:02Z",
            "endTime": "2025-07-09T14:03:03Z",
            "attributes": null,
            "droppedAttributesCount": 0,
            "events": null,
            "droppedEventsCount": 0,
            "links": null,
            "droppedLinksCount": 0,
            "status": null,
            "resource": null,
            "scope": null
        }
    ]
}
`,
		},
		{
			id:       "jq",
			trace:    &mackerel.TraceResponse{Spans: []*mackerel.Span{{Name: "test-span"}}},
			jqFilter: ".spans[].name",
			expected: `test-span
`,
		},
		{
			id:    "no spans",
			trace: &mackerel.TraceResponse{Spans: []*mackerel.Span{}},
			expected: `{
    "spans": []
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			client := mackerelclient.NewMockClient(
				mackerelclient.MockGetTrace(func(id string) (*mackerel.TraceResponse, error) {
					assert.Equal(t, traceID, id)
					return tc.trace, nil
				}),
			)
			out := new(bytes.Buffer)
			app := &tracesApp{
				client:    client,
				outStream: out,
				jqFilter:  tc.jqFilter,
			}
			assert.NoError(t, app.getTrace(t.Context(), traceID))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
