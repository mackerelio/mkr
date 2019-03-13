package services

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	mackerel "github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/mackerelclient"
)

func TestServicesApp_Run(t *testing.T) {
	testCases := []struct {
		id       string
		services []*mackerel.Service
		expected string
	}{
		{
			id: "default",
			services: []*mackerel.Service{
				&mackerel.Service{
					Name:  "sample-service-1",
					Memo:  "sample memo 1",
					Roles: []string{"role1", "role2", "role3"},
				},
				&mackerel.Service{
					Name:  "sample-service-2",
					Memo:  "sample memo 2",
					Roles: []string{"role"},
				},
				&mackerel.Service{
					Name:  "sample-service-3",
					Memo:  "",
					Roles: []string{},
				},
			},
			expected: `[
    {
        "name": "sample-service-1",
        "memo": "sample memo 1",
        "roles": [
            "role1",
            "role2",
            "role3"
        ]
    },
    {
        "name": "sample-service-2",
        "memo": "sample memo 2",
        "roles": [
            "role"
        ]
    },
    {
        "name": "sample-service-3",
        "memo": "",
        "roles": []
    }
]
`,
		},
		{
			id:       "no services",
			services: []*mackerel.Service{},
			expected: `[]
`,
		},
	}
	for _, tc := range testCases {
		client := mackerelclient.NewMockClient(
			mackerelclient.MockFindServices(func() ([]*mackerel.Service, error) {
				return tc.services, nil
			}),
		)
		t.Run(tc.id, func(t *testing.T) {
			out := new(bytes.Buffer)
			app := &servicesApp{
				client:    client,
				outStream: out,
			}
			assert.NoError(t, app.run())
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
