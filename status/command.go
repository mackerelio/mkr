package status

import (
	"os"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

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

	if argHostID == "" {
		if argHostID = mackerelclient.LoadHostIDFromConfig(confFile); argHostID == "" {
			cli.ShowCommandHelpAndExit(c, "status", 1)
		}
	}

	host, err := mackerelclient.NewFromContext(c).FindHost(argHostID)
	logger.DieIf(err)

	if isVerbose {
		err := format.PrettyPrintJSON(os.Stdout, host)
		logger.DieIf(err)
	} else {
		err := format.PrettyPrintJSON(os.Stdout, &format.Host{
			ID:            host.ID,
			Name:          host.Name,
			DisplayName:   host.DisplayName,
			Status:        host.Status,
			RoleFullnames: host.GetRoleFullnames(),
			IsRetired:     host.IsRetired,
			CreatedAt:     format.ISO8601Extended(host.DateFromCreatedAt()),
			IPAddresses:   host.IPAddresses(),
		})
		logger.DieIf(err)
	}
	return nil
}
