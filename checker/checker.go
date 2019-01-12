package checker

import (
	"github.com/mackerelio/mackerel-agent/config"
	cli "gopkg.in/urfave/cli.v1"
)

var Command = cli.Command{
	Name:  "run-checks",
	Usage: "run check commands in mackerel-agent.conf",
	Description: `
    Execute command of check plugins in mackerel-agent.conf all at once.
    It is used for checking setting and operation of the check plugins.
`,
	Action: doRunChecks,
}

func doRunChecks(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	conf, err := config.LoadConfig(confFile)
	if err != nil {
		return err
	}
	return runChecks(conf.CheckPlugins)
}

func runChecks(plugins map[string]*config.CheckPlugin) error {
	return nil
}
