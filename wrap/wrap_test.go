package wrap

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/urfave/cli/v3"
)

func newWrapCommand(t testing.TB, args []string) *cli.Command {
	t.Helper()
	var retVal cli.Command

	var cmd = *Command
	cmd.Action = func(_ context.Context, c *cli.Command) error {
		retVal = *c
		return nil
	}

	(&cli.Command{
		Commands: []*cli.Command{&cmd},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "conf",
			},
			&cli.StringFlag{
				Name: "apibase",
			},
		},
	}).Run(t.Context(), args) // nolint

	return &retVal
}

func TestCommand_Action(t *testing.T) {
	type testResult struct {
		Name                 string               `json:"name"`
		Status               mackerel.CheckStatus `json:"status"`
		Message              string               `json:"message"`
		NotificationInterval uint                 `json:"notificationInterval,omitempty"` // Minutes
	}
	type testReq struct {
		Reports []testResult `json:"reports"`
	}

	testCases := []struct {
		Name string
		Args []string

		Result   testResult
		ExitCode int
	}{
		{
			Name: "simple",
			Args: []string{
				"-name=test-check",
				"-detail",
				"-note", "This is note",
				"--",
				"go", "run", "testdata/stub.go",
			},
			Result: testResult{
				Name:   "test-check",
				Status: mackerel.CheckStatusCritical,
				Message: `command exited with code: 1
Note: This is note
% go run testdata/stub.go
Hello.
exit status 1
`,
				NotificationInterval: 0,
			},
			ExitCode: 1,
		},
		{
			Name: "long output",
			Args: []string{
				"-name=test-check",
				"-detail",
				"-note", "This is note",
				"--",
				"go", "run", "testdata/long.go",
			},
			Result: testResult{
				Name:   "test-check",
				Status: mackerel.CheckStatusCritical,
				Message: `command exited with code: 1
Note: This is note
% go run testdata/long.go
` + strings.Repeat("Hello world!\n", 33) + `Hello w
...
!
` + strings.Repeat("Hello world!\n", 38) + `exit status 1
`,
				NotificationInterval: 0,
			},
			ExitCode: 1,
		},
		{
			Name: "notification interval",
			Args: []string{
				"-name=test-check2",
				"-auto-close",
				"-notification-interval", "20m",
				"--",
				"echo", "1",
			},
			Result: testResult{
				Name:   "test-check2",
				Status: mackerel.CheckStatusOK,
				Message: `command exited with code: 0
% echo 1`,
				NotificationInterval: 20,
			},
			ExitCode: 0,
		},
		{
			Name: "minimum notification interval",
			Args: []string{
				"-name=test-check3",
				"-auto-close",
				"-notification-interval", "5m", // when less 10 min then 10 min.
				"--",
				"echo", "2",
			},
			Result: testResult{
				Name:   "test-check3",
				Status: mackerel.CheckStatusOK,
				Message: `command exited with code: 0
% echo 2`,
				NotificationInterval: 10,
			},
			ExitCode: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				reqPath := "/api/v0/monitoring/checks/report"
				if req.URL.Path != reqPath {
					t.Errorf("request URL should be %s but: %s", reqPath, req.URL.Path)
				}

				body, _ := io.ReadAll(req.Body)
				var treq testReq

				err := json.Unmarshal(body, &treq)
				if err != nil {
					t.Fatal("request body should be decoded as json", string(body))
				}
				got := treq.Reports[0]
				expect := tc.Result

				if !reflect.DeepEqual(got, expect) {
					t.Errorf("something went wrong.\n   got: %+v,\nexpect: %+v", got, expect)
				}

				res.Header()["Content-Type"] = []string{"application/json"}
				err = json.NewEncoder(res).Encode(map[string]bool{
					"success": true,
				})
				if err != nil {
					t.Fatal(err)
				}
			}))
			defer ts.Close()

			args := append(
				// This test checks to verify the sending request.
				// Therefore, to test in an environment without a configuration file, it needs to set the host argument.
				[]string{"$0", "-conf=testdata/dummy.conf", "-apibase", ts.URL, "wrap", "-host", "3Yr"},
				tc.Args...,
			)

			cmd := newWrapCommand(t, args)
			err := doWrap(t.Context(), cmd)
			var exitCode int
			if err != nil {
				exitCode = 1
				if excoder, ok := err.(cli.ExitCoder); ok {
					exitCode = excoder.ExitCode()
				}
			}
			if exitCode != tc.ExitCode {
				t.Errorf("exit code %d is expected. but: %d", tc.ExitCode, exitCode)
			}
		})
	}
}

// Commands are executed even if the configuration file does not exist.
func TestCommand_Action_withoutConf(t *testing.T) {
	cmd := newWrapCommand(t, []string{
		"$0", "-conf=notfound", "-apibase=http://localhost", "wrap",
		"--detail", "--",
		"go", "run", "testdata/stub.go",
	})
	expect := "command exited with code: 1"
	err := doWrap(t.Context(), cmd)
	if err == nil {
		t.Errorf("error should be occurred but nil")
	} else if err.Error() != expect {
		t.Errorf("The error message is different from the expected.\n   got: %s\nexpect: %s",
			err, expect)
	}
}

func Test_truncate(t *testing.T) {
	testCases := []struct {
		src      string
		limit    int
		sep      string
		expected string
	}{
		{
			src:      "",
			limit:    0,
			sep:      "",
			expected: "",
		},
		{
			src:      "",
			limit:    10,
			sep:      " ... ",
			expected: "",
		},
		{
			src:      "Hello, world!",
			limit:    100,
			sep:      " ... ",
			expected: "Hello, world!",
		},
		{
			src:      "Hello, world!",
			limit:    0,
			sep:      " ... ",
			expected: "",
		},
		{
			src:      "Hello, world!",
			limit:    3,
			sep:      " ... ",
			expected: " ..",
		},
		{
			src:      "Hello, world!",
			limit:    5,
			sep:      " ... ",
			expected: " ... ",
		},
		{
			src:      "Hello, world!",
			limit:    10,
			sep:      " ... ",
			expected: "He ... ld!",
		},
		{
			src:      "Hello, world!",
			limit:    15,
			sep:      " ... ",
			expected: "Hello, world!",
		},
		{
			src:      "こんにちは、世界",
			limit:    6,
			sep:      "..",
			expected: "こん..世界",
		},
		{
			src:      strings.Repeat("abcde", 10),
			limit:    30,
			sep:      " ... ",
			expected: "abcdeabcdeab ... cdeabcdeabcde",
		},
	}
	for _, tc := range testCases {
		got := truncate(tc.src, tc.limit, tc.sep)
		if got != tc.expected {
			t.Errorf("truncate(%q, %d, %q) should be %q but got: %q",
				tc.src, tc.limit, tc.sep, tc.expected, got)
		}
		if len([]rune(got)) > tc.limit {
			t.Errorf("length of truncate(%q, %d, %q) should not exceed %d but got: %d",
				tc.src, tc.limit, tc.sep, tc.limit, len([]rune(got)))
		}
	}
}
