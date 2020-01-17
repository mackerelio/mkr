package org

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/mackerelclient"
)

func TestOrgApp_Run(t *testing.T) {
	testCases := []struct {
		id       string
		org      *mackerel.Org
		expected string
	}{
		{
			id:  "default",
			org: &mackerel.Org{Name: "sample-org"},
			expected: `{
    "name": "sample-org"
}
`,
		},
	}
	for _, tc := range testCases {
		client := mackerelclient.NewMockClient(
			mackerelclient.MockGetOrg(func() (*mackerel.Org, error) {
				return tc.org, nil
			}),
		)
		t.Run(tc.id, func(t *testing.T) {
			out := new(bytes.Buffer)
			app := &orgApp{
				client:    client,
				outStream: out,
			}
			assert.NoError(t, app.run())
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
