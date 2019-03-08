package org

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

// Command is the definition of org subcommand
var Command = cli.Command{
	Name:  "org",
	Usage: "Fetch organization",
	Description: `
    Fetch organization.
    Requests APIs under "/api/v0/org". See https://mackerel.io/api-docs/entry/organizations .
`,
	Action: doOrg,
}

func doOrg(c *cli.Context) error {
	client := mackerelclient.NewFromContext(c)

	org, err := client.GetOrg()
	logger.DieIf(err)
	format.PrettyPrintJSON(os.Stdout, org)
	return nil
}
