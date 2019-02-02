package wrap

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mackerel "github.com/mackerelio/mackerel-client-go"
	cli "gopkg.in/urfave/cli.v1"
)

func newWrapContext(args []string) *cli.Context {
	app := cli.NewApp()
	parentFs := flag.NewFlagSet("mockmkr", flag.ContinueOnError)
	for _, f := range []cli.Flag{
		cli.StringFlag{Name: "conf"}, cli.StringFlag{Name: "apibase"},
	} {
		f.Apply(parentFs)
	}
	parentFs.Parse(args)
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
	fs.Parse(args)
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
		Name           string
		Args           []string
		ExpectedResult testResult
	}{
		{
			Name: "simple",
			Args: []string{
				"-name=test-check",
				"-detail",
				"-memo", "This is memo",
				"--",
				"go", "run", "testdata/stub.go",
			},
			ExpectedResult: testResult{
				Name:   "test-check",
				Status: mackerel.CheckStatusCritical,
				Message: `command exited with code: 1
Memo: This is memo
% go run testdata/stub.go
Hello.
exit status 1
`,
				NotificationInterval: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				reqPath := "/api/v0/monitoring/checks/report"
				if req.URL.Path != reqPath {
					t.Errorf("request URL should be %s but: %s", reqPath, req.URL.Path)
				}

				body, _ := ioutil.ReadAll(req.Body)
				var treq testReq

				err := json.Unmarshal(body, &treq)
				if err != nil {
					t.Fatal("request body should be decoded as json", string(body))
				}
				got := treq.Reports[0]
				expect := tc.ExpectedResult

				if !reflect.DeepEqual(got, expect) {
					t.Errorf("something went wrong.\n   got: %+v,\nexpect: %+v", got, expect)
				}

				res.Header()["Content-Type"] = []string{"application/json"}
				json.NewEncoder(res).Encode(map[string]bool{
					"success": true,
				})
			}))
			defer ts.Close()

			args := append(
				[]string{"-conf=testdata/dummy.conf", "-apibase", ts.URL, "wrap"},
				tc.Args...,
			)

			c := newWrapContext(args)
			Command.Action.(func(*cli.Context) error)(c)
		})
	}
}
