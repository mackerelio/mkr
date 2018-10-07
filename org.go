package main

import (
	"github.com/mackerelio/mkr/logger"
	"gopkg.in/urfave/cli.v1"
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
	client := newMackerelFromContext(c)

	org, err := client.GetOrg()
	logger.DieIf(err)
	PrettyPrintJSON(org)
	return nil
}
