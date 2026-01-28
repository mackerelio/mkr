package users

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/stretchr/testify/assert"

	"github.com/mackerelio/mkr/mackerelclient"
)

var (
	sampleUser1 = &mackerel.User{
		ID:                      "user1",
		ScreenName:              "Alice",
		Email:                   "alice@example.com",
		Authority:               "owner",
		IsInRegistrationProcess: false,
		IsMFAEnabled:            true,
		AuthenticationMethods:   []string{"password", "google"},
		JoinedAt:                1553000000,
	}
	sampleUser2 = &mackerel.User{
		ID:                      "user2",
		ScreenName:              "Bob",
		Email:                   "bob@example.com",
		Authority:               "collaborator",
		IsInRegistrationProcess: true,
		IsMFAEnabled:            false,
		AuthenticationMethods:   []string{"password"},
		JoinedAt:                1552000000,
	}
)

func TestUserApp_FindUsers(t *testing.T) {
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	defer func() { time.Local = nil }()
	testCases := []struct {
		id       string
		verbose  bool
		format   string
		users    []*mackerel.User
		expected string
	}{
		{
			id:    "default",
			users: []*mackerel.User{sampleUser1, sampleUser2},
			expected: `[
    {
        "id": "user1",
        "screenName": "Alice",
        "email": "alice@example.com",
        "authority": "owner",
        "isInRegistrationProcess": false,
        "isMFAEnabled": true,
        "authenticationMethods": [
            "password",
            "google"
        ],
        "joinedAt": "2019-03-19T21:53:20+09:00"
    },
    {
        "id": "user2",
        "screenName": "Bob",
        "email": "bob@example.com",
        "authority": "collaborator",
        "isInRegistrationProcess": true,
        "isMFAEnabled": false,
        "authenticationMethods": [
            "password"
        ],
        "joinedAt": "2019-03-08T08:06:40+09:00"
    }
]
`,
		},
		{
			id:      "verbose",
			users:   []*mackerel.User{sampleUser1, sampleUser2},
			verbose: true,
			expected: `[
    {
        "id": "user1",
        "screenName": "Alice",
        "email": "alice@example.com",
        "authority": "owner",
        "isMFAEnabled": true,
        "authenticationMethods": [
            "password",
            "google"
        ],
        "joinedAt": 1553000000
    },
    {
        "id": "user2",
        "screenName": "Bob",
        "email": "bob@example.com",
        "authority": "collaborator",
        "isInRegistrationProcess": true,
        "authenticationMethods": [
            "password"
        ],
        "joinedAt": 1552000000
    }
]
`,
		},
		{
			id:     "format",
			users:  []*mackerel.User{sampleUser1, sampleUser2},
			format: `{{range .}}{{.ID}} {{.ScreenName}} {{.Authority}} {{.JoinedAt}}{{"\n"}}{{end}}`,
			expected: `user1 Alice owner 1553000000
user2 Bob collaborator 1552000000
`,
		},
		{
			id:       "empty",
			users:    []*mackerel.User{},
			expected: "[]\n",
		},
	}
	for _, tc := range testCases {
		client := mackerelclient.NewMockClient(
			mackerelclient.MockFindUsers(func() ([]*mackerel.User, error) {
				return tc.users, nil
			}),
		)
		t.Run(tc.id, func(t *testing.T) {
			out := new(bytes.Buffer)
			app := &userApp{
				client:    client,
				outStream: out,
			}
			assert.NoError(t, app.findUsers(t.Context(), findUsersParam{
				verbose: tc.verbose,
				format:  tc.format,
			}))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}

func TestUserApp_FindUsersError(t *testing.T) {
	client := mackerelclient.NewMockClient(
		mackerelclient.MockFindUsers(func() ([]*mackerel.User, error) {
			return nil, fmt.Errorf("API error")
		}),
	)
	out := new(bytes.Buffer)
	app := &userApp{
		client:    client,
		outStream: out,
	}
	err := app.findUsers(t.Context(), findUsersParam{})
	assert.Error(t, err)
	assert.Equal(t, "API error", err.Error())
}

func TestUserApp_FindUsersFormatAndJqIncompatible(t *testing.T) {
	client := mackerelclient.NewMockClient(
		mackerelclient.MockFindUsers(func() ([]*mackerel.User, error) {
			return []*mackerel.User{sampleUser1}, nil
		}),
	)
	out := new(bytes.Buffer)
	app := &userApp{
		client:    client,
		outStream: out,
		jqFilter:  ".[]",
	}
	err := app.findUsers(t.Context(), findUsersParam{
		format: "{{.ID}}",
	})
	assert.Error(t, err)
	assert.Equal(t, "--format and --jq options are incompatible.", err.Error())
}
