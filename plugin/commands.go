package plugin

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

// CommandPlugin is definition of mkr plugin
var CommandPlugin = cli.Command{
	Name:        "plugin",
	Usage:       "Manage mackerel plugin",
	Description: `[WIP] Manage mackerel plugin`,
	Subcommands: []cli.Command{
		{
			Name:        "install",
			Usage:       "install mackerel plugin",
			Description: `WIP`,
			Action:      doPluginInstall,
		},
	},
	Hidden: true,
}

// main function for mkr plugin install
func doPluginInstall(c *cli.Context) error {
	fmt.Println("do plugin install")
	return nil
}
