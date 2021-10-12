package aws_integrations

import (
	"os"

	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:      "aws-integrations",
	Usage:     "List aws integration settings",
	ArgsUsage: "",
	Description: `
	List the information of the aws integration settings.
	Requests "GET /api/v0/aws-integrations". See https://mackerel.io/api-docs/entry/aws-integration#list.
`,
	Action: doAWSIntegrations,
}

func doAWSIntegrations(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}
	return (&awsIntegrationsApp{
		client:    client,
		outStream: os.Stdout,
	}).run()
}
