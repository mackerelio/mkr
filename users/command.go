package users

import (
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v2"
)

var CommandUsers = &cli.Command{
	Name:      "users",
	Usage:     "List users",
	ArgsUsage: "[--verbose | -v] [--format | -f <format>] [--jq <formula>]",
	Description: `
    List the information of the users.
    Requests "GET /api/v0/users". See https://mackerel.io/api-docs/entry/users#list .
`,
	Action: doUsers,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Value:   "",
			Usage:   "Output format template",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Verbose output mode",
		},
		jq.CommandLineFlag,
	},
}

func doUsers(c *cli.Context) error {
	client, err := mackerelclient.New(c.String("conf"), c.String("apibase"))
	if err != nil {
		return err
	}

	return (&userApp{
		client:    client,
		logger:    logger.New(),
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).findUsers(findUsersParam{
		verbose: c.Bool("verbose"),
		format:  c.String("format"),
	})
}
