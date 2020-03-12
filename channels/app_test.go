package channels

import (
	"bytes"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/stretchr/testify/assert"
)

// boolPointer is a helper function to initialize a bool pointer
func boolPointer(b bool) *bool {
	return &b
}

func TestChannelsApp_Run(t *testing.T) {
	testCases := []struct {
		id       string
		channels []*mackerel.Channel
		expected string
	}{
		{
			id: "default",
			channels: []*mackerel.Channel{
				&mackerel.Channel{
					ID:      "abcdefabc",
					Name:    "email channel",
					Type:    "email",
					Emails:  &[]string{"test@example.com", "test2@example.com"},
					UserIDs: &[]string{"1234", "2345"},
					Events:  &[]string{"alert"},
				},
				&mackerel.Channel{
					ID:   "bcdefabcd",
					Name: "slack channel",
					Type: "slack",
					URL:  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
					Mentions: mackerel.Mentions{
						OK:      "ok message",
						Warning: "warning message",
					},
					EnabledGraphImage: boolPointer(true),
					Events:            &[]string{"alert"},
				},
				&mackerel.Channel{
					ID:     "cdefabcde",
					Name:   "webhook channel",
					Type:   "webhook",
					URL:    "http://example.com/webhook",
					Events: &[]string{"alertGroup"},
				},
				&mackerel.Channel{
					ID:   "defabcdef",
					Name: "line channel",
					Type: "line",
				},
			},
			expected: `[
    {
        "id": "abcdefabc",
        "name": "email channel",
        "type": "email",
        "emails": [
            "test@example.com",
            "test2@example.com"
        ],
        "userIds": [
            "1234",
            "2345"
        ],
        "mentions": {},
        "events": [
            "alert"
        ]
    },
    {
        "id": "bcdefabcd",
        "name": "slack channel",
        "type": "slack",
        "mentions": {
            "ok": "ok message",
            "warning": "warning message"
        },
        "enabledGraphImage": true,
        "url": "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
        "events": [
            "alert"
        ]
    },
    {
        "id": "cdefabcde",
        "name": "webhook channel",
        "type": "webhook",
        "mentions": {},
        "url": "http://example.com/webhook",
        "events": [
            "alertGroup"
        ]
    },
    {
        "id": "defabcdef",
        "name": "line channel",
        "type": "line",
        "mentions": {}
    }
]
`,
		},
		{
			id:       "no channels",
			channels: []*mackerel.Channel{},
			expected: `[]
`,
		},
	}
	for _, tc := range testCases {
		client := mackerelclient.NewMockClient(
			mackerelclient.MockFindChannels(func() ([]*mackerel.Channel, error) {
				return tc.channels, nil
			}),
		)
		t.Run(tc.id, func(t *testing.T) {
			out := new(bytes.Buffer)
			app := &channelsApp{
				client:    client,
				outStream: out,
			}
			assert.NoError(t, app.run())
			assert.Equal(t, tc.expected, out.String())
		})
	}
}

// TODO: write tests for pullChannels
