package status

import (
	"fmt"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/mackerelclient"
	cli "gopkg.in/urfave/cli.v1"
)

// Command is definition of mkr status subcommand
var Command = cli.Command{
	Name:      "status",
	Usage:     "Show the host",
	ArgsUsage: "[--verbose | -v] <hostId>",
	Description: `
    Show the information of the host identified with <hostId>.
    Requests "GET /api/v0/hosts/<hostId>". See https://mackerel.io/api-docs/entry/hosts#get .
`,
	Action: doStatus,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
	},
}

func doStatus(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	argHostID := c.Args().Get(0)
	isVerbose := c.Bool("verbose")

	mcli, err := mackerelclient.New(confFile, c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	if argHostID == "" {
		if argHostID = mackerelclient.LoadHostIDFromConfig(confFile); argHostID == "" {
			return fmt.Errorf("no hostIDs are specified from args or idfile")
		}
	}
	host, err := mcli.FindHost(argHostID)
	if err != nil {
		return err
	}

	var stuff interface{}
	if isVerbose {
		stuff = host
	} else {
		stuff = &format.Host{
			ID:            host.ID,
			Name:          host.Name,
			DisplayName:   host.DisplayName,
			Status:        host.Status,
			RoleFullnames: host.GetRoleFullnames(),
			IsRetired:     host.IsRetired,
			CreatedAt:     format.ISO8601Extended(host.DateFromCreatedAt()),
			IPAddresses:   host.IPAddresses(),
		}
	}
	return format.PrettyPrintJSON(stuff)
}
