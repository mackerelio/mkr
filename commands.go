package main

import (
	"strings"

	"github.com/Songmu/prompter"
	"github.com/mackerelio/mkr/aws_integrations"
	"github.com/mackerelio/mkr/channels"
	"github.com/mackerelio/mkr/checks"
	"github.com/mackerelio/mkr/fetch"
	"github.com/mackerelio/mkr/hosts"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/mackerelio/mkr/metrics"
	"github.com/mackerelio/mkr/org"
	"github.com/mackerelio/mkr/plugin"
	"github.com/mackerelio/mkr/services"
	"github.com/mackerelio/mkr/status"
	"github.com/mackerelio/mkr/throw"
	"github.com/mackerelio/mkr/update"
	"github.com/mackerelio/mkr/wrap"
	"github.com/urfave/cli"
)

// Commands cli.Command object list
var Commands = []cli.Command{
	status.Command,
	hosts.CommandHosts,
	hosts.CommandCreate,
	update.Command,
	throw.Command,
	metrics.Command,
	fetch.Command,
	commandRetire,
	services.Command,
	commandMonitors,
	channels.Command,
	commandAlerts,
	commandDashboards,
	commandAnnotations,
	org.Command,
	plugin.CommandPlugin,
	checks.Command,
	wrap.Command,
	aws_integrations.Command,
}

var commandRetire = cli.Command{
	Name:      "retire",
	Usage:     "Retire hosts",
	ArgsUsage: "[--force] hostIds...",
	Description: `
    Retire host identified by <hostId>. Be careful because this is an irreversible operation.
    Requests POST /api/v0/hosts/<hostId>/retire parallelly. See https://mackerel.io/api-docs/entry/hosts#retire .
`,
	Action: doRetire,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "force", Usage: "Force retirement without confirmation."},
	},
}

func doRetire(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	force := c.Bool("force")
	argHostIDs := c.Args()

	if len(argHostIDs) < 1 {
		argHostIDs = make([]string, 1)
		if argHostIDs[0] = mackerelclient.LoadHostIDFromConfig(confFile); argHostIDs[0] == "" {
			cli.ShowCommandHelpAndExit(c, "retire", 1)
		}
	}

	if !force && !prompter.YN("Retire following hosts.\n  "+strings.Join(argHostIDs, "\n  ")+"\nAre you sure?", true) {
		logger.Log("", "retirement is canceled.")
		return nil
	}

	client := mackerelclient.NewFromContext(c)

	for _, hostID := range argHostIDs {
		err := client.RetireHost(hostID)
		logger.DieIf(err)

		logger.Log("retired", hostID)
	}
	return nil
}
