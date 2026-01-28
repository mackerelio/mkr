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
		jqFilter string
		expected string
	}{
		{
			id:       "default",
			org:      &mackerel.Org{Name: "sample-org"},
			jqFilter: "",
			expected: `{
    "name": "sample-org"
}
`,
		},
		{
			id:       "jq_orgName",
			org:      &mackerel.Org{Name: "sample-org"},
			jqFilter: ".name",
			expected: `sample-org
`,
		},
		{
			id:       "jq_emptyDisplayName",
			org:      &mackerel.Org{Name: "sample-org"},
			jqFilter: ".displayName",
			expected: "\n",
		},
		{
			id:       "jq_displayName",
			org:      &mackerel.Org{Name: "sample-org", DisplayName: "Sample Org"},
			jqFilter: ".displayName",
			expected: `Sample Org
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
			jqFilter := tc.jqFilter
			app := &orgApp{
				client:    client,
				outStream: out,
				jqFilter:  jqFilter,
			}
			assert.NoError(t, app.run(t.Context()))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
