package traces

import (
	"context"
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v3"
)

// Command is the definition of traces subcommand
var Command = &cli.Command{
	Name:  "traces",
	Usage: "Fetch trace information",
	Description: `
    Fetch trace information. With "get" subcommand, get detailed trace information for the specified trace ID.
    Requests APIs under "/api/v0/traces". See https://mackerel.io/api-docs/entry/traces .
`,
	Commands: []*cli.Command{
		{
			Name:      "get",
			Usage:     "get trace",
			ArgsUsage: "<traceId> [--jq <formula>]",
			Description: `
    Get detailed trace information for the specified trace ID.
    Requests "GET /api/v0/traces/<traceId>". See https://mackerel.io/api-docs/entry/traces#get .
`,
			Action: doTracesGet,
			Flags: []cli.Flag{
				jq.CommandLineFlag,
			},
		},
	},
}

func doTracesGet(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	traceID := c.Args().Get(0)
	client, err := mackerelclient.New(c.String("conf"), c.String("apibase"))
	if err != nil {
		return err
	}

	return (&tracesApp{
		client:    client,
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).getTrace(ctx, traceID)
}
