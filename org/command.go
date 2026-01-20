package org

import (
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v2"
)

// Command is the definition of org subcommand
var Command = cli.Command{
	Name:      "org",
	Usage:     "Fetch organization",
	ArgsUsage: "[--jq <formula>]",
	Description: `
    Fetch organization.
    Requests APIs under "/api/v0/org". See https://mackerel.io/api-docs/entry/organizations .
`,
	Action: doOrg,
	Flags: []cli.Flag{
		jq.CommandLineFlag,
	},
}

func doOrg(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&orgApp{
		client:    client,
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).run()
}
