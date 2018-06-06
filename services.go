package main

import (
	"os"

	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"gopkg.in/urfave/cli.v1"
)

var commandServices = cli.Command{
	Name:      "services",
	Usage:     "List services",
	ArgsUsage: "",
	Description: `
    Manipulate services. With no subcommand specified, this will list the information of the services.
    Requests APIs under "/api/v0/services". See https://mackerel.io/api-docs/entry/services .
`,
	Action: doServicesList,
	Flags:  []cli.Flag{},
	Subcommands: []cli.Command{
		{
			Name:      "create",
			Usage:     "create a new service",
			ArgsUsage: "[--memo | -m <memo>] serviceName",
			Description: `
    Create a new service with given name.
    Requests "POST /api/v0/services". See https://mackerel.io/api-docs/entry/services#create .
`,
			Action: doServiceCreate,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "memo, m", Value: "", Usage: "Memo for the service"},
			},
		},
	},
}

func doServicesList(c *cli.Context) error {
	services, err := newMackerelFromContext(c).FindServices()
	logger.DieIf(err)
	PrettyPrintJSON(services)
	return nil
}

func doServiceCreate(c *cli.Context) error {
	argServiceName := c.Args().Get(0)
	optMemo := c.String("memo")

	if argServiceName == "" {
		cli.ShowSubcommandHelp(c)
		os.Exit(1)
	}

	client := newMackerelFromContext(c)

	service, err := client.CreateService(&mkr.CreateServiceParam{
		Name: argServiceName,
		Memo: optMemo,
	})
	logger.DieIf(err)

	logger.Log("created", service.Name)

	return nil
}
