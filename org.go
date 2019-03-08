package main

import (
	"os"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	cli "gopkg.in/urfave/cli.v1"
)

var commandOrg = cli.Command{
	Name:  "org",
	Usage: "Fetch organization",
	Description: `
    Fetch organization.
    Requests APIs under "/api/v0/org". See https://mackerel.io/api-docs/entry/organizations .
`,
	Action: doOrgRetrieve,
}

func doOrgRetrieve(c *cli.Context) error {
	client := mackerelclient.NewFromContext(c)

	org, err := client.GetOrg()
	logger.DieIf(err)
	format.PrettyPrintJSON(os.Stdout, org)
	return nil
}
