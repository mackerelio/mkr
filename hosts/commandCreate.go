package hosts

import (
	"os"

	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

// CommandCreate is definition of mkr create subcommand
var CommandCreate = cli.Command{
	Name:      "create",
	Usage:     "Create a new host",
	ArgsUsage: "[--status | -st <status>] [--roleFullname | -R <service:role>] [--customIdentifier <customIdentifier>] <hostName>",
	Description: `
    Create a new host with status, roleFullname and/or customIdentifier.
    Requests "POST /api/v0/hosts". See https://mackerel.io/api-docs/entry/hosts#create .
`,
	Action: doCreate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "status, st", Value: "", Usage: "Host status ('working', 'standby', 'maintenance')"},
		cli.StringSliceFlag{
			Name:  "roleFullname, R",
			Value: &cli.StringSlice{},
			Usage: "Multiple choices are allowed. ex. My-Service:proxy, My-Service:db-master",
		},
		cli.StringFlag{Name: "customIdentifier", Value: "", Usage: "CustomIdentifier for the Host"},
	},
}

func doCreate(c *cli.Context) error {
	argHostName := c.Args().Get(0)
	if argHostName == "" {
		cli.ShowCommandHelp(c, "create")
		os.Exit(1)
	}

	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&hostApp{
		client:    client,
		logger:    logger.New(),
		outStream: os.Stdout,
	}).createHost(createHostParam{
		name:             argHostName,
		roleFullnames:    c.StringSlice("roleFullname"),
		status:           c.String("status"),
		customIdentifier: c.String("customIdentifier"),
	})
}
