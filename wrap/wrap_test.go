package wrap

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/urfave/cli"
)

func newWrapContext(t testing.TB, args []string) *cli.Context {
	t.Helper()
	app := cli.NewApp()
	parentFs := flag.NewFlagSet("mockmkr", flag.ContinueOnError)
	for _, f := range []cli.Flag{
		cli.StringFlag{Name: "conf"}, cli.StringFlag{Name: "apibase"},
	} {
		f.Apply(parentFs)
	}
	if err := parentFs.Parse(args); err != nil {
		t.Fatal(err)
	}
	for i, v := range parentFs.Args() {
		if v == "wrap" {
			args = parentFs.Args()[i+1:]
			break
		}
	}
	parentCtx := cli.NewContext(app, parentFs, nil)

	fs := flag.NewFlagSet("mockwrap", flag.ContinueOnError)
	for _, f := range Command.Flags {
		f.Apply(fs)
	}
	if err := fs.Parse(args); err != nil {
		t.Fatal(err)
	}
	return cli.NewContext(app, fs, parentCtx)
}

func TestCommand_Action(t *testing.T) {
	type testResult struct {
		Name                 string               `json:"name"`
		Status               mackerel.CheckStatus `json:"status"`
		Message              string               `json:"message"`
		NotificationInterval uint                 `json:"notificationInterval,omitempty"`
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
% echo 1
`,
				NotificationInterval: 1200,
			},
			ExitCode: 0,
		},
		{
			Name: "minimum notification interval",
			Args: []string{
				"-name=test-check3",
				"-auto-close",
				"-notification-interval", "5m",
				"--",
				"echo", "2",
			},
			Result: testResult{
				Name:   "test-check2",
				Status: mackerel.CheckStatusOK,
				Message: `command exited with code: 0
% echo 2
`,
				NotificationInterval: 600,
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
				[]string{"-conf=testdata/dummy.conf", "-apibase", ts.URL, "wrap"},
				tc.Args...,
			)

			c := newWrapContext(t, args)
			err := Command.Action.(func(*cli.Context) error)(c)
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

func TestCommand_Action_withoutConf(t *testing.T) {
	c := newWrapContext(t, []string{
		"-conf=notfound", "-apibase=http://localhost", "wrap",
		"--detail", "--",
		"go", "run", "testdata/stub.go",
	})
	expect := "command exited with code: 1"
	err := Command.Action.(func(*cli.Context) error)(c)
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
