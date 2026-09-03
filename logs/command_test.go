package logs

import (
	"context"
	"strings"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/urfave/cli/v3"
)

// runBuildFindLogsParam parses args through the "find" subcommand's flags
// and returns whatever buildFindLogsParam produces, without touching the
// network (the Action is replaced so mackerelclient.New is never called).
func runBuildFindLogsParam(t *testing.T, args []string) (*mackerel.FindLogsParam, error) {
	t.Helper()

	var gotParam *mackerel.FindLogsParam
	var gotErr error

	cmd := &cli.Command{
		Name:  "find",
		Flags: findLogsFlags(),
		Action: func(ctx context.Context, c *cli.Command) error {
			gotParam, gotErr = buildFindLogsParam(c)
			return nil
		},
	}

	if err := cmd.Run(t.Context(), append([]string{"find"}, args...)); err != nil {
		return nil, err
	}
	return gotParam, gotErr
}

func TestBuildFindLogsParam(t *testing.T) {
	baseArgs := []string{"--service", "shoppingcart", "--from", "1700000000", "--to", "1700003600"}

	tests := map[string]struct {
		args    []string
		wantErr string
		check   func(t *testing.T, param *mackerel.FindLogsParam)
	}{
		"basic": {
			args: baseArgs,
			check: func(t *testing.T, param *mackerel.FindLogsParam) {
				if param.ServiceName != "shoppingcart" {
					t.Errorf("ServiceName should be shoppingcart but: %v", param.ServiceName)
				}
				if param.From.Unix() != 1700000000 {
					t.Errorf("From should be 1700000000 but: %v", param.From.Unix())
				}
				if param.To.Unix() != 1700003600 {
					t.Errorf("To should be 1700003600 but: %v", param.To.Unix())
				}
				if param.First != nil {
					t.Errorf("First should be nil but: %v", *param.First)
				}
				if param.Last != nil {
					t.Errorf("Last should be nil but: %v", *param.Last)
				}
			},
		},
		"first_positive": {
			args: append(append([]string{}, baseArgs...), "--first", "50"),
			check: func(t *testing.T, param *mackerel.FindLogsParam) {
				if param.First == nil || *param.First != 50 {
					t.Errorf("First should be 50 but: %v", param.First)
				}
			},
		},
		"first_zero_is_error": {
			args:    append(append([]string{}, baseArgs...), "--first", "0"),
			wantErr: "--first must be positive",
		},
		"first_negative_is_error": {
			args:    append(append([]string{}, baseArgs...), "--first", "-1"),
			wantErr: "--first must be positive",
		},
		"last_positive": {
			args: append(append([]string{}, baseArgs...), "--last", "20"),
			check: func(t *testing.T, param *mackerel.FindLogsParam) {
				if param.Last == nil || *param.Last != 20 {
					t.Errorf("Last should be 20 but: %v", param.Last)
				}
			},
		},
		"last_zero_is_error": {
			args:    append(append([]string{}, baseArgs...), "--last", "0"),
			wantErr: "--last must be positive",
		},
		"last_negative_is_error": {
			args:    append(append([]string{}, baseArgs...), "--last", "-5"),
			wantErr: "--last must be positive",
		},
		"severities_and_keywords": {
			args: append(append([]string{}, baseArgs...), "--severity", "ERROR", "--severity", "FATAL", "--keyword", "timeout"),
			check: func(t *testing.T, param *mackerel.FindLogsParam) {
				wantSeverities := []mackerel.LogSeverity{mackerel.LogSeverityError, mackerel.LogSeverityFatal}
				if len(param.Severities) != len(wantSeverities) {
					t.Fatalf("Severities should be %v but: %v", wantSeverities, param.Severities)
				}
				for i, s := range wantSeverities {
					if param.Severities[i] != s {
						t.Errorf("Severities[%d] should be %v but: %v", i, s, param.Severities[i])
					}
				}
				if len(param.Keywords) != 1 || param.Keywords[0] != "timeout" {
					t.Errorf("Keywords should be [timeout] but: %v", param.Keywords)
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			param, err := runBuildFindLogsParam(t, tt.args)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q but got: %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tt.check(t, param)
		})
	}
}
