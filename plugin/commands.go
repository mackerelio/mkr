package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	argInstallTarget := c.Args().First()
	if argInstallTarget == "" {
		return fmt.Errorf("Specify install name")
	}

	_, err := parseInstallTarget(argInstallTarget)
	if err != nil {
		return errors.Wrap(err, "failed to install plugin")
	}

	err = setupPluginDir(c.String("prefix"))
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

type installTarget struct {
	owner      string
	repo       string
	pluginName string
	releaseTag string
}

// Parse install target string passed from args
// example is below
// - mackerelio/mackerel-plugin-sample
// - mackerel-plugin-sample
// - mackerelio/mackerel-plugin-sample@v0.0.1
func parseInstallTarget(target string) (*installTarget, error) {
	it := &installTarget{}

	ownerRepoAndReleaseTag := strings.Split(target, "@")
	var ownerRepo string
	switch len(ownerRepoAndReleaseTag) {
	case 1:
		ownerRepo = ownerRepoAndReleaseTag[0]
	case 2:
		ownerRepo = ownerRepoAndReleaseTag[0]
		it.releaseTag = ownerRepoAndReleaseTag[1]
	default:
		return nil, fmt.Errorf("Install target is invalid: %s", target)
	}

	ownerAndRepo := strings.Split(ownerRepo, "/")
	switch len(ownerAndRepo) {
	case 1:
		it.pluginName = ownerAndRepo[0]
	case 2:
		it.owner = ownerAndRepo[0]
		it.repo = ownerAndRepo[1]
	default:
		return nil, fmt.Errorf("Install target is invalid: %s", target)
	}

	return it, nil
}
