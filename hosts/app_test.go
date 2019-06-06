package hosts

import (
	"bytes"
	"testing"
	"time"

	mackerel "github.com/mackerelio/mackerel-client-go"
	"github.com/stretchr/testify/assert"

	"github.com/mackerelio/mkr/mackerelclient"
)

var (
	sampleHost1 = &mackerel.Host{
		ID:          "foo",
		Name:        "sample.app1",
		DisplayName: "Sample Host foo",
		Status:      mackerel.HostStatusWorking,
		Roles: mackerel.Roles{
			"SampleService": []string{"app"},
		},
		IsRetired: false,
		CreatedAt: 1553000000,
		Interfaces: []mackerel.Interface{
			{
				Name:      "en0",
				IPAddress: "10.0.0.1",
			},
		},
	}
	sampleHost2 = &mackerel.Host{
		ID:          "bar",
		Name:        "sample.app2",
		DisplayName: "Sample Host bar",
		Status:      mackerel.HostStatusStandby,
		Roles: mackerel.Roles{
			"SampleService": []string{"db"},
		},
		IsRetired: false,
		CreatedAt: 1552000000,
		Interfaces: []mackerel.Interface{
			{
				Name:      "eth0",
				IPAddress: "10.0.1.2",
			},
		},
	}
)

func TestHostApp_FindHosts(t *testing.T) {
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	defer func() { time.Local = nil }()
	testCases := []struct {
		id       string
		verbose  bool
		name     string
		service  string
		roles    []string
		statuses []string
		format   string
		hosts    []*mackerel.Host
		expected string
	}{
		{
			id:    "default",
			hosts: []*mackerel.Host{sampleHost1, sampleHost2},
			expected: `[
    {
        "id": "foo",
        "name": "sample.app1",
        "displayName": "Sample Host foo",
        "status": "working",
        "roleFullnames": [
            "SampleService:app"
        ],
        "isRetired": false,
        "createdAt": "2019-03-19T21:53:20+09:00",
        "ipAddresses": {
            "en0": "10.0.0.1"
        }
    },
    {
        "id": "bar",
        "name": "sample.app2",
        "displayName": "Sample Host bar",
        "status": "standby",
        "roleFullnames": [
            "SampleService:db"
        ],
        "isRetired": false,
        "createdAt": "2019-03-08T08:06:40+09:00",
        "ipAddresses": {
            "eth0": "10.0.1.2"
        }
    }
]
`,
		},
		{
			id:      "verbose",
			hosts:   []*mackerel.Host{sampleHost1, sampleHost2},
			verbose: true,
			expected: `[
    {
        "id": "foo",
        "name": "sample.app1",
        "displayName": "Sample Host foo",
        "type": "",
        "status": "working",
        "memo": "",
        "roles": {
            "SampleService": [
                "app"
            ]
        },
        "isRetired": false,
        "createdAt": 1553000000,
        "meta": {},
        "interfaces": [
            {
                "name": "en0",
                "ipAddress": "10.0.0.1"
            }
        ]
    },
    {
        "id": "bar",
        "name": "sample.app2",
        "displayName": "Sample Host bar",
        "type": "",
        "status": "standby",
        "memo": "",
        "roles": {
            "SampleService": [
                "db"
            ]
        },
        "isRetired": false,
        "createdAt": 1552000000,
        "meta": {},
        "interfaces": [
            {
                "name": "eth0",
                "ipAddress": "10.0.1.2"
            }
        ]
    }
]
`,
		},
		{
			id:     "format",
			hosts:  []*mackerel.Host{sampleHost1, sampleHost2},
			format: `{{range .}}{{.ID}} {{.Name}} {{.Status}} {{.CreatedAt}}{{"\n"}}{{end}}`,
			expected: `foo sample.app1 working 1553000000
bar sample.app2 standby 1552000000
`,
		},
		{
			id:       "name",
			hosts:    []*mackerel.Host{},
			name:     "Sample.app",
			expected: "null\n",
		},
		{
			id:       "service",
			hosts:    []*mackerel.Host{},
			service:  "SampleService",
			expected: "null\n",
		},
		{
			id:       "roles",
			hosts:    []*mackerel.Host{},
			roles:    []string{"role1", "role2"},
			expected: "null\n",
		},
		{
			id:       "statuses",
			hosts:    []*mackerel.Host{},
			statuses: []string{mackerel.HostStatusPoweroff, mackerel.HostStatusMaintenance},
			expected: "null\n",
		},
	}
	for _, tc := range testCases {
		client := mackerelclient.NewMockClient(
			mackerelclient.MockFindHosts(func(param *mackerel.FindHostsParam) ([]*mackerel.Host, error) {
				assert.Equal(t, tc.name, param.Name)
				assert.Equal(t, tc.service, param.Service)
				assert.Equal(t, tc.roles, param.Roles)
				assert.Equal(t, tc.statuses, param.Statuses)
				return tc.hosts, nil
			}),
		)
		t.Run(tc.id, func(t *testing.T) {
			out := new(bytes.Buffer)
			app := &hostApp{
				client:    client,
				outStream: out,
			}
			assert.NoError(t, app.findHosts(findHostsParam{
				verbose:  tc.verbose,
				name:     tc.name,
				service:  tc.service,
				roles:    tc.roles,
				statuses: tc.statuses,
				format:   tc.format,
			}))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
