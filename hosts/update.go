package hosts

import (
	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v2"
)

var CommandUpdate = cli.Command{
	Name:      "update",
	Usage:     "Update the host",
	ArgsUsage: "[--name | -n <name>] [--displayName <displayName>] [--status | -st <status>] [--roleFullname | -R <service:role>] [--overwriteRoles | -o] [--memo <memo>] [<hostIds...>]",
	Description: `
    Update the host identified with <hostId>.
    Requests "PUT /api/v0/hosts/<hostId>". See https://mackerel.io/api-docs/entry/hosts#update-information .
`,
	Action: doUpdate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Value: "",
			Usage: "Update hostname.",
		},
		cli.StringFlag{
			Name:  "displayName",
			Value: "",
			Usage: "Update displayName.",
		},
		cli.StringFlag{
			Name:  "status, st",
			Value: "",
			Usage: "Update status.",
		},
		cli.StringSliceFlag{
			Name:  "roleFullname, R",
			Value: &cli.StringSlice{},
			Usage: "Update rolefullname.",
		},
		cli.BoolFlag{
			Name:  "overwriteRoles, o",
			Usage: "Overwrite roles instead of adding specified roles.",
		},
		cli.StringFlag{
			Name:  "memo",
			Value: "",
			Usage: "memo for the Host",
		},
	},
}

func doUpdate(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	argHostIDs := c.Args()
	optName := c.String("name")
	optDisplayName := c.String("displayName")
	optStatus := c.String("status")
	optRoleFullnames := c.StringSlice("roleFullname")
	overwriteRoles := c.Bool("overwriteRoles")
	optMemo := c.String("memo")

	if len(argHostIDs) < 1 {
		argHostIDs = make([]string, 1)
		if argHostIDs[0] = mackerelclient.LoadHostIDFromConfig(confFile); argHostIDs[0] == "" {
			cli.ShowCommandHelpAndExit(c, "update", 1)
		}
	}

	needUpdateHostStatus := optStatus != ""
	needUpdateRolesInHostUpdate := !overwriteRoles && len(optRoleFullnames) > 0
	needUpdateHost := (optName != "" || optDisplayName != "" || overwriteRoles || optMemo != "" || needUpdateRolesInHostUpdate)

	if !needUpdateHostStatus && !needUpdateHost {
		logger.Log("update", "at least one argumet is required.")
		cli.ShowCommandHelpAndExit(c, "update", 1)
	}

	client := mackerelclient.NewFromContext(c)

	for _, hostID := range argHostIDs {
		if needUpdateHostStatus {
			err := client.UpdateHostStatus(hostID, optStatus)
			logger.DieIf(err)
		}

		if overwriteRoles {
			err := client.UpdateHostRoleFullnames(hostID, optRoleFullnames)
			logger.DieIf(err)
		}

		if needUpdateHost {
			host, err := client.FindHost(hostID)
			logger.DieIf(err)
			name := ""
			if optName == "" {
				name = host.Name
			} else {
				name = optName
			}
			displayname := ""
			if optDisplayName == "" {
				displayname = host.DisplayName
			} else {
				displayname = optDisplayName
			}
			memo := ""
			if optMemo == "" {
				memo = host.Memo
			} else {
				memo = optMemo
			}
			param := &mackerel.UpdateHostParam{
				Name:        name,
				DisplayName: displayname,
				Meta:        host.Meta,
				Interfaces:  host.Interfaces,
				Memo:        memo,
			}
			if needUpdateRolesInHostUpdate {
				param.RoleFullnames = optRoleFullnames
			}
			_, err = client.UpdateHost(hostID, param)
			logger.DieIf(err)
		}

		logger.Log("updated", hostID)
	}
	return nil
}
