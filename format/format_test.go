package format

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestISO8601Extended(t *testing.T) {
	now := time.Now()
	expect := now.Format("2006-01-02T15:04:05-07:00") // ISO 8601 extended format
	got := ISO8601Extended(now)
	if got != expect {
		t.Errorf("should be %q got %q", expect, got)
	}
}

func TestPrettyPrintJSON(t *testing.T) {
	testCases := []struct {
		id       string
		srcJson  string
		query    string
		expected string
	}{
		{
			id:      "simple",
			srcJson: `{"id": 1, "string": "foo"}`,
			expected: `{
    "id": 1,
    "string": "foo"
}
`,
		},
		{
			id:       "simple with query",
			srcJson:  `{"id": 1, "string": "foo"}`,
			expected: "1\n",
			query:    ".id",
		},
		{
			id:      "array",
			srcJson: `[{"id": 1, "string": "foo"},{"id": 2, "string": "bar"}]`,
			expected: `[
    {
        "id": 1,
        "string": "foo"
    },
    {
        "id": 2,
        "string": "bar"
    }
]
`,
		},
		{
			id:       "array with query",
			srcJson:  `[{"id": 1, "string": "foo"},{"id": 2, "string": "bar"}]`,
			expected: "1\n2\n",
			query:    ".[] | .id",
		},
		{
			id:      "array with select object",
			srcJson: `[{"id": 1, "string": "foo"},{"id": 2, "string": "bar"}]`,
			expected: `{"id":2,"string":"bar"}
`,
			query: ".[] | select(.id == 2)",
		},
		{
			id:       "array with select value",
			srcJson:  `[{"id": 1, "string": "foo"},{"id": 2, "string": "bar"}]`,
			expected: "bar\n",
			query:    ".[] | select(.id == 2) | .string",
		},
		{
			id:      "simple with brackets",
			srcJson: `{"id": 1, "string": "\u003cfoo\u003e"}`,
			expected: `{
    "id": 1,
    "string": "<foo>"
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			out := new(bytes.Buffer)
			var src interface{}
			assert.NoError(t, json.Unmarshal([]byte(tc.srcJson), &src))
			assert.NoError(t, PrettyPrintJSON(out, src, tc.query))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
