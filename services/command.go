package services

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

// Command is the definition of services subcommand
var Command = cli.Command{
	Name:      "services",
	Usage:     "List services",
	ArgsUsage: "",
	Description: `
    List the information of the services.
    Requests "GET /api/v0/services". See https://mackerel.io/api-docs/entry/services#list.
`,
	Action: doServices,
}

func doServices(c *cli.Context) error {
	services, err := mackerelclient.NewFromContext(c).FindServices()
	logger.DieIf(err)
	format.PrettyPrintJSON(os.Stdout, services)
	return nil
}
