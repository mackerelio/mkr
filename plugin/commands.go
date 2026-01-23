package plugin

import (
	"github.com/urfave/cli/v2"
)

// CommandPlugin is definition of mkr plugin
var CommandPlugin = &cli.Command{
	Name:  "plugin",
	Usage: "Manage mackerel plugin",
	Description: `
    Manage mackerel plugin.  For example, you can install a mackerel plugin and
    check plugin by "mkr plugin install".
`,
	Subcommands: []*cli.Command{
		commandPluginInstall,
	},
}
