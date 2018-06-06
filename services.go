package main

import (
	"github.com/mackerelio/mkr/logger"
	"gopkg.in/urfave/cli.v1"
)

var commandServices = cli.Command{
	Name:      "services",
	Usage:     "List services",
	ArgsUsage: "",
	Description: `
    List the information of the services.
    Requests "GET /api/v0/services". See https://mackerel.io/api-docs/entry/services#list.
`,
	Action: doServices,
	Flags:  []cli.Flag{},
}

func doServices(c *cli.Context) error {
	services, err := newMackerelFromContext(c).FindServices()
	logger.DieIf(err)
	PrettyPrintJSON(services)
	return nil
}
