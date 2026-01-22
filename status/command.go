package status

import (
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "status",
	Usage:     "Show the host",
	ArgsUsage: "[--verbose | -v] [--jq <formula>] <hostId>",
	Description: `
    Show the information of the host identified with <hostId>.
    Requests "GET /api/v0/hosts/<hostId>". See https://mackerel.io/api-docs/entry/hosts#get .
`,
	Action: doStatus,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Verbose output mode",
		},
		jq.CommandLineFlag,
	},
}

func doStatus(c *cli.Context) error {
	confFile := c.String("conf")
	argHostID := c.Args().Get(0)
	isVerbose := c.Bool("verbose")

	if argHostID == "" {
		if argHostID = mackerelclient.LoadHostIDFromConfig(confFile); argHostID == "" {
			cli.ShowCommandHelpAndExit(c, "status", 1)
		}
	}

	client := mackerelclient.NewFromContext(c)

	return (&statussApp{
		client:    client,
		outStream: os.Stdout,
		isVerbose: isVerbose,
		argHostID: argHostID,

		jqFilter: c.String("jq"),
	}).run()
}
