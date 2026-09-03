package logs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v3"
)

// Command is the definition of logs subcommand
var Command = &cli.Command{
	Name:  "logs",
	Usage: "Search logs",
	Description: `
    Search logs. Requests APIs under "/api/v0/logs". See https://mackerel.io/api-docs/entry/logs .
`,
	Commands: []*cli.Command{
		{
			Name:  "find",
			Usage: "find logs",
			ArgsUsage: "--service <service> --from <from> --to <to> " +
				"[--service-namespace <namespace>] [--keyword <keyword>] [--severity <severity>] [--trace-id <traceId>] " +
				"[--first <first> --after <after> | --last <last> --before <before>] [--jq <formula>]",
			Description: `
    Search logs matching the specified conditions.
    Requests "POST /api/v0/logs". See https://mackerel.io/api-docs/entry/logs#find .

    --first/--after (forward pagination) and --last/--before (backward pagination) are mutually
    exclusive: --first cannot be combined with --last or --before, and --last cannot be combined
    with --first or --after.
`,
			Action: doLogsFind,
			Flags:  findLogsFlags(),
		},
	},
}

// findLogsFlags returns the flag definitions for "logs find".
// cli.Flag instances retain internal state (hasBeenSet etc.) across Run calls,
// so a fresh set is created each time instead of being reused (mainly relevant for tests).
func findLogsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "service",
			Usage:    "Service name.",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "service-namespace",
			Usage: "Service namespace.",
		},
		&cli.Int64Flag{
			Name:     "from",
			Usage:    "Search from (unix timestamp).",
			Required: true,
		},
		&cli.Int64Flag{
			Name:     "to",
			Usage:    "Search to (unix timestamp).",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "keyword",
			Usage: "Keyword to search in log body. Can be specified multiple times (AND condition).",
		},
		&cli.StringSliceFlag{
			Name:  "severity",
			Usage: "Severity to filter (UNSPECIFIED, TRACE, DEBUG, INFO, WARN, ERROR, FATAL). Can be specified multiple times.",
		},
		&cli.StringFlag{
			Name:  "trace-id",
			Usage: "Trace ID to filter.",
		},
		&cli.IntFlag{
			Name:  "first",
			Usage: "Number of logs to fetch from the beginning. Must be positive. Cannot be used with --last or --before.",
		},
		&cli.StringFlag{
			Name:  "after",
			Usage: "Cursor to fetch logs after (cannot be used with --last or --before).",
		},
		&cli.IntFlag{
			Name:  "last",
			Usage: "Number of logs to fetch from the end. Must be positive. Cannot be used with --first or --after.",
		},
		&cli.StringFlag{
			Name:  "before",
			Usage: "Cursor to fetch logs before (cannot be used with --first or --after).",
		},
		jq.CommandLineFlag,
	}
}

func doLogsFind(ctx context.Context, c *cli.Command) error {
	client, err := mackerelclient.New(c.String("conf"), c.String("apibase"))
	if err != nil {
		return err
	}

	param, err := buildFindLogsParam(c)
	if err != nil {
		return err
	}

	return (&logsApp{
		client:    client,
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).findLogs(ctx, param)
}

func buildFindLogsParam(c *cli.Command) (*mackerel.FindLogsParam, error) {
	param := &mackerel.FindLogsParam{
		ServiceName: c.String("service"),
		From:        time.Unix(c.Int64("from"), 0),
		To:          time.Unix(c.Int64("to"), 0),
		Keywords:    c.StringSlice("keyword"),
	}

	if ns := c.String("service-namespace"); ns != "" {
		param.ServiceNamespace = &ns
	}
	if traceID := c.String("trace-id"); traceID != "" {
		param.TraceID = &traceID
	}
	if after := c.String("after"); after != "" {
		param.After = &after
	}
	if before := c.String("before"); before != "" {
		param.Before = &before
	}
	if c.IsSet("first") {
		first := c.Int("first")
		if first <= 0 {
			return nil, fmt.Errorf("--first must be positive, but got %d", first)
		}
		param.First = &first
	}
	if c.IsSet("last") {
		last := c.Int("last")
		if last <= 0 {
			return nil, fmt.Errorf("--last must be positive, but got %d", last)
		}
		param.Last = &last
	}

	for _, s := range c.StringSlice("severity") {
		severity := mackerel.LogSeverity(s)
		param.Severities = append(param.Severities, severity)
	}

	return param, nil
}
