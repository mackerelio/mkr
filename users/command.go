package users

import (
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var CommandUsers = cli.Command{
	Name:      "users",
	Usage:     "List users",
	ArgsUsage: "[--verbose | -v] [--format | -f <format>] [--jq <formula>]",
	Description: `
    List the information of the users.
    Requests "GET /api/v0/users". See https://mackerel.io/ja/api-docs/entry/users#list .
`,
	Action: doUsers,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "format, f", Value: "", Usage: "Output format template"},
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
		jq.CommandLineFlag,
	},
}

func doUsers(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
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
