package traces

import (
	"bytes"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/stretchr/testify/assert"
)

func TestTracesApp_getTrace(t *testing.T) {
	traceID := "0123456789abcdef0123456789abcdef"
	
	// Sample trace response
	expectedTrace := &mackerel.TraceResponse{
		Spans: []*mackerel.Span{
			{
				TraceID:   traceID,
				SpanID:    "012345678012345678",
				Name:      "test-span",
				Kind:      "internal",
				StartTime: "2025-07-09T14:03:02.000Z",
				EndTime:   "2025-07-09T14:03:02.000Z",
			},
		},
	}

	client := mackerelclient.NewMockClient(
		mackerelclient.MockGetTrace(func(id string) (*mackerel.TraceResponse, error) {
			assert.Equal(t, traceID, id)
			return expectedTrace, nil
		}),
	)

	out := new(bytes.Buffer)
	app := &tracesApp{
		client:    client,
		outStream: out,
		jqFilter:  "",
	}

	err := app.getTrace(traceID)
	assert.NoError(t, err)
	
	// Check that JSON output contains expected data
	output := out.String()
	assert.Contains(t, output, traceID)
	assert.Contains(t, output, "test-span")
	assert.Contains(t, output, "internal")
}