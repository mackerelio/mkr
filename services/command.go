package services

import (
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

// Command is the definition of services subcommand
var Command = cli.Command{
	Name:      "services",
	Usage:     "List services",
	ArgsUsage: "[--jq <formula>]",
	Description: `
    List the information of the services.
    Requests "GET /api/v0/services". See https://mackerel.io/api-docs/entry/services#list.
`,
	Action: doServices,
	Flags: []cli.Flag{
		jq.CommandLineFlag,
	},
}

func doServices(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&servicesApp{
		client:    client,
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).run()
}
