package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v3"
)

const (
	maxPeriod  = 30 * 24 * time.Hour
	maxPerPage = 1000
)

var traceIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)

var severities = []mackerel.LogSeverity{
	mackerel.LogSeverityUNSPECIFIED,
	mackerel.LogSeverityTRACE,
	mackerel.LogSeverityDEBUG,
	mackerel.LogSeverityINFO,
	mackerel.LogSeverityWARN,
	mackerel.LogSeverityERROR,
	mackerel.LogSeverityFATAL,
}

var allowedOperators = map[mackerel.LogBodyAttributeValueType][]mackerel.LogBodyAttributeOperator{
	mackerel.LogBodyAttributeValueTypeString: {
		mackerel.LogBodyAttributeOperatorEQ,
		mackerel.LogBodyAttributeOperatorNEQ,
		mackerel.LogBodyAttributeOperatorSTARTSWITH,
	},
	mackerel.LogBodyAttributeValueTypeInt: {
		mackerel.LogBodyAttributeOperatorEQ,
		mackerel.LogBodyAttributeOperatorGT,
		mackerel.LogBodyAttributeOperatorGTE,
		mackerel.LogBodyAttributeOperatorLT,
		mackerel.LogBodyAttributeOperatorLTE,
	},
	mackerel.LogBodyAttributeValueTypeDouble: {
		mackerel.LogBodyAttributeOperatorEQ,
		mackerel.LogBodyAttributeOperatorGT,
		mackerel.LogBodyAttributeOperatorGTE,
		mackerel.LogBodyAttributeOperatorLT,
		mackerel.LogBodyAttributeOperatorLTE,
	},
	mackerel.LogBodyAttributeValueTypeBool: {
		mackerel.LogBodyAttributeOperatorEQ,
		mackerel.LogBodyAttributeOperatorNEQ,
	},
}

// Command is the definition of logs subcommand
var Command = &cli.Command{
	Name:  "logs",
	Usage: "Search logs and list saved log searches",
	Description: `
    Search logs and list saved log searches.
    Requests APIs under "/api/v0/logs" and "/api/v0/saved-log-searches". See https://mackerel.io/api-docs/entry/logs .
`,
	Commands: []*cli.Command{
		{
			Name:  "search",
			Usage: "search logs",
			ArgsUsage: "--service|-s <service> [--namespace <namespace>] --from <from> [--to <to>] [--keyword|-k <keyword>]... [--severity <severity>]... " +
				"[--trace-id <traceId>] [--body-attribute <json>]... [--order <ASC|DESC>] [--per-page <n>] [--cursor <cursor> [--cursor-direction <NEXT|PREVIOUS>]] [--jq <formula>]",
			Description: `
    Search logs of the specified service. The time range is [from, to) and must not exceed 30 days.
    If --to is omitted, the current time is used.
    The response contains "pagination.startCursor" and "pagination.endCursor";
    pass one of them to --cursor to get the next (or previous) page.
    Requests "POST /api/v0/logs". See https://mackerel.io/api-docs/entry/logs#list .

    --body-attribute filters on a field of a JSON log body and takes a JSON object, e.g.
      --body-attribute '{"keyPath":["req","method"],"operator":"EQ","type":"string","value":"GET"}'
      --body-attribute '{"keyPath":["status"],"operator":"GTE","type":"int","value":"500"}'
    keyPath, value, operator and type are all required.
    value is always a JSON string, whatever the type is (e.g. "500", "1.5", "true").
    operator: EQ, NEQ, GT, GTE, LT, LTE, STARTS_WITH
    type: string, int, double, bool
    Specify --severity or --keyword multiple times to give several values. Values are not split by commas.
`,
			Action:                    doLogsSearch,
			DisableSliceFlagSeparator: true,
			Flags:                     searchFlags(),
		},
		{
			Name:      "saved-searches",
			Usage:     "list saved log searches",
			ArgsUsage: "[--jq <formula>]",
			Description: `
    List saved log searches.
    Requests "GET /api/v0/saved-log-searches". See https://mackerel.io/api-docs/entry/logs#list-saved-searches .
`,
			Action: doLogsSavedSearches,
			Flags: []cli.Flag{
				jq.CommandLineFlag,
			},
		},
	},
}

func searchFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "Service name (required)",
		},
		&cli.StringFlag{
			Name:  "namespace",
			Usage: "Service namespace",
		},
		&cli.Int64Flag{
			Name:  "from",
			Usage: "The start of the period to search, inclusive (epoch seconds, required)",
		},
		&cli.Int64Flag{
			Name:  "to",
			Usage: "The end of the period to search, exclusive (epoch seconds). Defaults to the current time",
		},
		&cli.StringSliceFlag{
			Name:    "keyword",
			Aliases: []string{"k"},
			Usage:   "Substring that the log body must contain. If specified multiple times, all must be contained",
		},
		&cli.StringSliceFlag{
			Name:  "severity",
			Usage: "Severity of logs to match (UNSPECIFIED, TRACE, DEBUG, INFO, WARN, ERROR, FATAL). If specified multiple times, any may match",
		},
		&cli.StringFlag{
			Name:  "trace-id",
			Usage: "Trace ID (32 hexadecimal characters)",
		},
		&cli.StringSliceFlag{
			Name:  "body-attribute",
			Usage: `Condition on a field of a JSON log body, as a JSON object like '{"keyPath":["req","method"],"operator":"EQ","type":"string","value":"GET"}'. value is always a string, even for int, double and bool. All conditions must match`,
		},
		&cli.StringFlag{
			Name:  "order",
			Usage: "Sort direction by timestamp (ASC or DESC). Defaults to DESC",
		},
		&cli.IntFlag{
			Name:  "per-page",
			Usage: "Number of logs per page (1-1000). Defaults to 50",
		},
		&cli.StringFlag{
			Name:  "cursor",
			Usage: "Cursor to start paging from (pagination.startCursor or pagination.endCursor of a previous response). The log at the cursor is excluded",
		},
		&cli.StringFlag{
			Name:  "cursor-direction",
			Usage: "Direction to read from --cursor (NEXT or PREVIOUS). Defaults to NEXT",
		},
		jq.CommandLineFlag,
	}
}

func doLogsSearch(ctx context.Context, c *cli.Command) error {
	param, err := newSearchParam(c)
	if err != nil {
		return err
	}
	client, err := mackerelclient.New(c.String("conf"), c.String("apibase"))
	if err != nil {
		return err
	}
	return (&logsApp{
		client:    client,
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).search(ctx, param)
}

func doLogsSavedSearches(ctx context.Context, c *cli.Command) error {
	client, err := mackerelclient.New(c.String("conf"), c.String("apibase"))
	if err != nil {
		return err
	}
	return (&logsApp{
		client:    client,
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).listSavedSearches(ctx)
}

func newSearchParam(c *cli.Command) (*mackerel.ListLogsParam, error) {
	service := c.String("service")
	if service == "" {
		return nil, cli.Exit("`service` is a required field to search logs.", 1)
	}
	if !c.IsSet("from") {
		return nil, cli.Exit("`from` is a required field to search logs.", 1)
	}
	from := c.Int64("from")
	to := time.Now().Unix()
	if c.IsSet("to") {
		to = c.Int64("to")
	}
	if to < from {
		return nil, cli.Exit("`to` must not be earlier than `from`.", 1)
	}
	if to-from > int64(maxPeriod/time.Second) {
		return nil, cli.Exit("the period between `from` and `to` must be within 30 days.", 1)
	}

	param := &mackerel.ListLogsParam{ServiceName: service, From: from, To: to}

	if c.IsSet("namespace") {
		v := c.String("namespace")
		param.ServiceNamespace = &v
	}
	if c.IsSet("keyword") {
		for _, k := range c.StringSlice("keyword") {
			if k == "" {
				return nil, cli.Exit("`keyword` must not be empty.", 1)
			}
		}
		param.Keywords = c.StringSlice("keyword")
	}
	if c.IsSet("severity") {
		for _, s := range c.StringSlice("severity") {
			sev := mackerel.LogSeverity(strings.ToUpper(s))
			if !slices.Contains(severities, sev) {
				return nil, cli.Exit(fmt.Sprintf("invalid `severity`: %q (must be one of %v).", s, severities), 1)
			}
			param.Severities = append(param.Severities, sev)
		}
	}
	if c.IsSet("trace-id") {
		v := c.String("trace-id")
		if !traceIDPattern.MatchString(v) {
			return nil, cli.Exit("`trace-id` must be 32 hexadecimal characters.", 1)
		}
		param.TraceID = &v
	}
	if c.IsSet("body-attribute") {
		for _, s := range c.StringSlice("body-attribute") {
			f, err := parseBodyAttribute(s)
			if err != nil {
				return nil, err
			}
			param.BodyAttributes = append(param.BodyAttributes, *f)
		}
	}
	if c.IsSet("order") {
		dir := mackerel.OrderDirection(strings.ToUpper(c.String("order")))
		if dir != mackerel.OrderDirectionASC && dir != mackerel.OrderDirectionDESC {
			return nil, cli.Exit(fmt.Sprintf("invalid `order`: %q (must be ASC or DESC).", c.String("order")), 1)
		}
		param.Order = &mackerel.LogOrder{Direction: &dir}
	}
	if c.IsSet("per-page") {
		v := c.Int("per-page")
		if v < 1 || v > maxPerPage {
			return nil, cli.Exit("`per-page` must be between 1 and 1000.", 1)
		}
		param.PerPage = &v
	}
	if c.IsSet("cursor-direction") && !c.IsSet("cursor") {
		return nil, cli.Exit("`cursor-direction` requires `cursor`.", 1)
	}
	if c.IsSet("cursor") {
		cursor := c.String("cursor")
		if cursor == "" {
			return nil, cli.Exit("`cursor` must not be empty.", 1)
		}
		dir := mackerel.LogPageCursorDirectionNEXT
		if c.IsSet("cursor-direction") {
			dir = mackerel.LogPageCursorDirection(strings.ToUpper(c.String("cursor-direction")))
			if dir != mackerel.LogPageCursorDirectionNEXT && dir != mackerel.LogPageCursorDirectionPREVIOUS {
				return nil, cli.Exit(fmt.Sprintf("invalid `cursor-direction`: %q (must be NEXT or PREVIOUS).", c.String("cursor-direction")), 1)
			}
		}
		param.PageCursor = &mackerel.LogPageCursor{Direction: dir, Cursor: cursor}
	}
	return param, nil
}

func parseBodyAttribute(s string) (*mackerel.LogBodyAttributeFilter, error) {
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	dec.DisallowUnknownFields()
	// Value is a pointer to tell an omitted value from an empty string.
	var v struct {
		KeyPath  []string                           `json:"keyPath"`
		Value    *string                            `json:"value"`
		Operator mackerel.LogBodyAttributeOperator  `json:"operator"`
		Type     mackerel.LogBodyAttributeValueType `json:"type"`
	}
	if err := dec.Decode(&v); err != nil {
		return nil, cli.Exit(fmt.Sprintf("invalid `body-attribute` %q: %v", s, err), 1)
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, cli.Exit(fmt.Sprintf("invalid `body-attribute` %q: unexpected data after the JSON object.", s), 1)
	}
	if len(v.KeyPath) == 0 {
		return nil, cli.Exit(fmt.Sprintf("invalid `body-attribute` %q: `keyPath` is required.", s), 1)
	}
	if v.Value == nil {
		return nil, cli.Exit(fmt.Sprintf("invalid `body-attribute` %q: `value` is required.", s), 1)
	}
	f := mackerel.LogBodyAttributeFilter{KeyPath: v.KeyPath, Value: *v.Value, Operator: v.Operator, Type: v.Type}
	ops, ok := allowedOperators[f.Type]
	if !ok {
		return nil, cli.Exit(fmt.Sprintf("invalid `body-attribute` %q: `type` must be one of string, int, double, bool.", s), 1)
	}
	if !slices.Contains(ops, f.Operator) {
		return nil, cli.Exit(fmt.Sprintf("invalid `body-attribute` %q: `operator` %q is not available for type %q (available: %v).", s, f.Operator, f.Type, ops), 1)
	}
	return &f, nil
}
