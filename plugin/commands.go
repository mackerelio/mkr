package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
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
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix",
					Usage: "plugin install location",
				},
			},
		},
	},
	Hidden: true,
}

// main function for mkr plugin install
func doPluginInstall(c *cli.Context) error {
	err := setupPluginDir(c.String("prefix"))
	if err != nil {
		return errors.Wrap(err, "failed to install plugin")
	}
	fmt.Println("do plugin install [wip]")
	return nil
}

// Create a directory for plugin install
func setupPluginDir(prefix string) error {
	if prefix == "" {
		prefix = "/opt/mackerel-agent/plugins"
	}
	err := os.MkdirAll(filepath.Join(prefix, "bin"), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to setup plugin directory")
	}
	return nil
}
