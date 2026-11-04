package logs

import (
	"bytes"
	"context"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestLogsApp_search(t *testing.T) {
	param := &mackerel.ListLogsParam{ServiceName: "shop", From: 100, To: 200}

	testCases := []struct {
		id       string
		res      *mackerel.ListLogsResponse
		jqFilter string
		expected string
	}{
		{
			id: "default",
			res: &mackerel.ListLogsResponse{
				Results: []*mackerel.LogEntry{{
					Cursor:      "c1",
					Severity:    mackerel.LogSeverityERROR,
					Body:        "failed",
					ServiceName: "shop",
				}},
				Pagination: &mackerel.LogPagination{HasNextPage: true, EndCursor: mackerel.ToPtr("c1")},
			},
			expected: `{
    "results": [
        {
            "cursor": "c1",
            "timestamp": "0001-01-01T00:00:00Z",
            "severity": "ERROR",
            "severityNumber": 0,
            "body": "failed",
            "traceId": null,
            "spanId": null,
            "serviceName": "shop",
            "serviceNamespace": null,
            "attributes": null,
            "resourceAttributes": null,
            "scopeAttributes": null
        }
    ],
    "pagination": {
        "hasNextPage": true,
        "hasPreviousPage": false,
        "startCursor": null,
        "endCursor": "c1"
    }
}
`,
		},
		{
			id: "jq",
			res: &mackerel.ListLogsResponse{
				Results: []*mackerel.LogEntry{{Body: "failed"}},
			},
			jqFilter: ".results[].body",
			expected: `failed
`,
		},
		{
			id:  "no logs",
			res: &mackerel.ListLogsResponse{Results: []*mackerel.LogEntry{}},
			expected: `{
    "results": [],
    "pagination": null
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			client := mackerelclient.NewMockClient(
				mackerelclient.MockListLogs(func(p *mackerel.ListLogsParam) (*mackerel.ListLogsResponse, error) {
					assert.Equal(t, param, p)
					return tc.res, nil
				}),
			)
			out := new(bytes.Buffer)
			app := &logsApp{client: client, outStream: out, jqFilter: tc.jqFilter}
			assert.NoError(t, app.search(t.Context(), param))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}

func TestLogsApp_listSavedSearches(t *testing.T) {
	testCases := []struct {
		id       string
		res      *mackerel.FindSavedLogSearchesResponse
		jqFilter string
		expected string
	}{
		{
			id: "default",
			res: &mackerel.FindSavedLogSearchesResponse{
				Results: []*mackerel.SavedLogSearch{{
					ID:   "s1",
					Name: "errors",
					Definition: &mackerel.SavedLogSearchDefinition{
						ServiceName: "shop",
						Severities:  []mackerel.LogSeverity{mackerel.LogSeverityERROR},
					},
					CreatedAt: 1700000000,
					UpdatedAt: 1700000100,
				}},
			},
			expected: `{
    "results": [
        {
            "id": "s1",
            "name": "errors",
            "description": null,
            "definition": {
                "serviceName": "shop",
                "serviceNamespace": null,
                "keywords": null,
                "severities": [
                    "ERROR"
                ],
                "traceId": null,
                "bodyAttributes": null
            },
            "createdAt": 1700000000,
            "updatedAt": 1700000100
        }
    ]
}
`,
		},
		{
			id:       "jq",
			res:      &mackerel.FindSavedLogSearchesResponse{Results: []*mackerel.SavedLogSearch{{Name: "errors"}}},
			jqFilter: ".results[].name",
			expected: `errors
`,
		},
		{
			id:  "no saved searches",
			res: &mackerel.FindSavedLogSearchesResponse{Results: []*mackerel.SavedLogSearch{}},
			expected: `{
    "results": []
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			client := mackerelclient.NewMockClient(
				mackerelclient.MockFindSavedLogSearches(func() (*mackerel.FindSavedLogSearchesResponse, error) {
					return tc.res, nil
				}),
			)
			out := new(bytes.Buffer)
			app := &logsApp{client: client, outStream: out, jqFilter: tc.jqFilter}
			assert.NoError(t, app.listSavedSearches(t.Context()))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}

func TestLogsApp_search_apiError(t *testing.T) {
	client := mackerelclient.NewMockClient(
		mackerelclient.MockListLogs(func(param *mackerel.ListLogsParam) (*mackerel.ListLogsResponse, error) {
			return nil, &mackerel.APIError{StatusCode: 403, Message: "log feature is not available"}
		}),
	)
	out := new(bytes.Buffer)
	app := &logsApp{client: client, outStream: out}
	assert.EqualError(t, app.search(t.Context(), &mackerel.ListLogsParam{ServiceName: "shop"}), "API request failed: log feature is not available")
	assert.Empty(t, out.String())
}

func TestNewSearchParam(t *testing.T) {
	testCases := []struct {
		id          string
		args        []string
		expected    *mackerel.ListLogsParam
		expectedErr string
	}{
		{
			id:   "required only",
			args: []string{"--service", "shop", "--from", "100", "--to", "200"},
			expected: &mackerel.ListLogsParam{
				ServiceName: "shop", From: 100, To: 200,
			},
		},
		{
			id: "all flags",
			args: []string{
				"-s", "shop", "--namespace", "ns", "--from", "100", "--to", "200",
				"-k", "a,b", "-k", "c",
				"--severity", "warn", "--severity", "ERROR",
				"--trace-id", "0123456789abcdef0123456789abcdef",
				"--body-attribute", `{"keyPath":["req","method"],"operator":"EQ","type":"string","value":"GET,POST"}`,
				"--body-attribute", `{"keyPath":["status"],"operator":"GTE","type":"int","value":"500"}`,
				"--order", "asc", "--per-page", "10",
				"--cursor", "c1", "--cursor-direction", "previous",
			},
			expected: &mackerel.ListLogsParam{
				ServiceName:      "shop",
				ServiceNamespace: mackerel.ToPtr("ns"),
				From:             100,
				To:               200,
				Keywords:         []string{"a,b", "c"},
				Severities:       []mackerel.LogSeverity{mackerel.LogSeverityWARN, mackerel.LogSeverityERROR},
				TraceID:          mackerel.ToPtr("0123456789abcdef0123456789abcdef"),
				BodyAttributes: []mackerel.LogBodyAttributeFilter{
					{KeyPath: []string{"req", "method"}, Operator: mackerel.LogBodyAttributeOperatorEQ, Type: mackerel.LogBodyAttributeValueTypeString, Value: "GET,POST"},
					{KeyPath: []string{"status"}, Operator: mackerel.LogBodyAttributeOperatorGTE, Type: mackerel.LogBodyAttributeValueTypeInt, Value: "500"},
				},
				Order:      &mackerel.LogOrder{Direction: mackerel.ToPtr(mackerel.OrderDirectionASC)},
				PageCursor: &mackerel.LogPageCursor{Direction: mackerel.LogPageCursorDirectionPREVIOUS, Cursor: "c1"},
				PerPage:    mackerel.ToPtr(10),
			},
		},
		{
			id:   "cursor defaults to NEXT",
			args: []string{"-s", "shop", "--from", "100", "--to", "200", "--cursor", "c1"},
			expected: &mackerel.ListLogsParam{
				ServiceName: "shop", From: 100, To: 200,
				PageCursor: &mackerel.LogPageCursor{Direction: mackerel.LogPageCursorDirectionNEXT, Cursor: "c1"},
			},
		},
		{
			id:          "no service",
			args:        []string{"--from", "100", "--to", "200"},
			expectedErr: "`service` is a required field to search logs.",
		},
		{
			id:          "no from",
			args:        []string{"-s", "shop", "--to", "200"},
			expectedErr: "`from` is a required field to search logs.",
		},
		{
			id:          "to before from",
			args:        []string{"-s", "shop", "--from", "200", "--to", "100"},
			expectedErr: "`to` must not be earlier than `from`.",
		},
		{
			id:          "period too long",
			args:        []string{"-s", "shop", "--from", "0", "--to", "2592001"},
			expectedErr: "the period between `from` and `to` must be within 30 days.",
		},
		{
			id:          "empty keyword",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "-k", ""},
			expectedErr: "`keyword` must not be empty.",
		},
		{
			id:          "invalid severity",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--severity", "WARN,ERROR"},
			expectedErr: `invalid ` + "`severity`" + `: "WARN,ERROR" (must be one of [UNSPECIFIED TRACE DEBUG INFO WARN ERROR FATAL]).`,
		},
		{
			id:          "invalid trace id",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--trace-id", "abc"},
			expectedErr: "`trace-id` must be 32 hexadecimal characters.",
		},
		{
			id:          "invalid order",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--order", "up"},
			expectedErr: `invalid ` + "`order`" + `: "up" (must be ASC or DESC).`,
		},
		{
			id:          "per-page out of range",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--per-page", "1001"},
			expectedErr: "`per-page` must be between 1 and 1000.",
		},
		{
			id:          "cursor-direction without cursor",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--cursor-direction", "NEXT"},
			expectedErr: "`cursor-direction` requires `cursor`.",
		},
		{
			id:          "invalid cursor-direction",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--cursor", "c", "--cursor-direction", "x"},
			expectedErr: `invalid ` + "`cursor-direction`" + `: "x" (must be NEXT or PREVIOUS).`,
		},
		{
			id:          "body-attribute not json",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--body-attribute", "foo"},
			expectedErr: "invalid `body-attribute` \"foo\": invalid character 'o' in literal false (expecting 'a')",
		},
		{
			id:          "body-attribute unknown field",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--body-attribute", `{"key":"a"}`},
			expectedErr: "invalid `body-attribute` \"{\\\"key\\\":\\\"a\\\"}\": json: unknown field \"key\"",
		},
		{
			id:          "body-attribute without keyPath",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--body-attribute", `{"operator":"EQ","type":"string","value":"a"}`},
			expectedErr: "invalid `body-attribute` \"{\\\"operator\\\":\\\"EQ\\\",\\\"type\\\":\\\"string\\\",\\\"value\\\":\\\"a\\\"}\": `keyPath` is required.",
		},
		{
			id:          "body-attribute without value",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--body-attribute", `{"keyPath":["a"],"operator":"EQ","type":"string"}`},
			expectedErr: "invalid `body-attribute` \"{\\\"keyPath\\\":[\\\"a\\\"],\\\"operator\\\":\\\"EQ\\\",\\\"type\\\":\\\"string\\\"}\": `value` is required.",
		},
		{
			id:          "body-attribute with trailing data",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--body-attribute", `{"keyPath":["a"],"operator":"EQ","type":"string","value":"x"} {}`},
			expectedErr: "invalid `body-attribute` \"{\\\"keyPath\\\":[\\\"a\\\"],\\\"operator\\\":\\\"EQ\\\",\\\"type\\\":\\\"string\\\",\\\"value\\\":\\\"x\\\"} {}\": unexpected data after the JSON object.",
		},
		{
			id:          "body-attribute with trailing closing brace",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--body-attribute", `{"keyPath":["a"],"operator":"EQ","type":"string","value":"x"}}`},
			expectedErr: "invalid `body-attribute` \"{\\\"keyPath\\\":[\\\"a\\\"],\\\"operator\\\":\\\"EQ\\\",\\\"type\\\":\\\"string\\\",\\\"value\\\":\\\"x\\\"}}\": unexpected data after the JSON object.",
		},
		{
			id:          "body-attribute invalid type",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--body-attribute", `{"keyPath":["a"],"operator":"EQ","type":"x","value":"a"}`},
			expectedErr: "invalid `body-attribute` \"{\\\"keyPath\\\":[\\\"a\\\"],\\\"operator\\\":\\\"EQ\\\",\\\"type\\\":\\\"x\\\",\\\"value\\\":\\\"a\\\"}\": `type` must be one of string, int, double, bool.",
		},
		{
			id:          "body-attribute operator not allowed for type",
			args:        []string{"-s", "shop", "--from", "100", "--to", "200", "--body-attribute", `{"keyPath":["a"],"operator":"GT","type":"string","value":"a"}`},
			expectedErr: "invalid `body-attribute` \"{\\\"keyPath\\\":[\\\"a\\\"],\\\"operator\\\":\\\"GT\\\",\\\"type\\\":\\\"string\\\",\\\"value\\\":\\\"a\\\"}\": `operator` \"GT\" is not available for type \"string\" (available: [EQ NEQ STARTS_WITH]).",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			var (
				param *mackerel.ListLogsParam
				err   error
			)
			cmd := *Command.Command("search")
			cmd.Flags = searchFlags() // flags keep their state across runs, so use fresh ones for each case
			cmd.Action = func(_ context.Context, c *cli.Command) error {
				param, err = newSearchParam(c)
				return nil
			}
			assert.NoError(t, cmd.Run(t.Context(), append([]string{"search"}, tc.args...)))
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, param)
		})
	}
}
