package services

import (
	"os"

	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
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
	Flags: []cli.Flag{
		cli.StringFlag{Name: "jq", Usage: "Query to select values from the response using jq syntax"},
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
		jq:        c.String("jq"),
	}).run()
}
