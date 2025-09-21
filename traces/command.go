package traces

import (
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

// Command is the definition of traces subcommand
var Command = cli.Command{
	Name:  "traces",
	Usage: "Fetch trace information",
	Description: `
    Fetch trace information. With "get" subcommand, get detailed trace information for the specified trace ID.
    Requests APIs under "/api/v0/traces". See https://mackerel.io/api-docs/entry/traces .
`,
	Subcommands: []cli.Command{
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

func doTracesGet(c *cli.Context) error {
	if len(c.Args()) != 1 {
		cli.ShowCommandHelpAndExit(c, "get", 1)
	}

	traceID := c.Args().Get(0)
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&tracesApp{
		client:    client,
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).getTrace(traceID)
}
