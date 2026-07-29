package apm

import (
	"bytes"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mackerelio/mkr/mackerelclient"
)

func TestDoHTTPServerStats(t *testing.T) {
	client := mackerelclient.NewMockClient(mackerelclient.MockListHTTPServerStats(func(param *mackerel.ListHTTPServerStatsParam) (*mackerel.HTTPServerStatsPageConnection, error) {
		assert.Equal(t, "service-name", param.ServiceName)
		return &mackerel.HTTPServerStatsPageConnection{}, nil
	}))
	var buf bytes.Buffer
	app := &httpServerStatsApp{
		client:    client,
		logger:    &testLogger{&buf},
		outStream: &buf,
	}
	err := app.listHTTPServerStats(t.Context(), &mackerel.ListHTTPServerStatsParam{
		ServiceName: "service-name",
	})
	require.NoError(t, err)
}
