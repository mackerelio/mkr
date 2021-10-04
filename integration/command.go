package integration

import (
	"os"

	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var CommandIntegration = cli.Command{
	Name:      "integration",
	Usage:     "List cloud integration settings",
	ArgsUsage: "[--provider | -p] <provider>",
	Description: `
	List the information of the cloud integration settings.
	For AWS, Requests "GET /api/v0/aws-integrations". See https://mackerel.io/api-docs/entry/aws-integration#list.
`,
	Action: doIntegration,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "provider, p", Usage: "view specific cloud integration settings"},
	},
}

func doIntegration(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}
	return (&integrationApp{
		client:    client,
		outStream: os.Stdout,
	}).run(
		c.String("provider"),
	)
}
