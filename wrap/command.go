package wrap

import (
	"fmt"
	"os"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mkr/logger"
	cli "gopkg.in/urfave/cli.v1"
)

// Command is definition of mkr wrap
var Command = cli.Command{
	Name:      "wrap",
	Usage:     "wrap command status",
	ArgsUsage: "[--verbose | -v] [--name | -n <name>] [--memo | -m <memo>] -- /path/to/batch",
	Description: `
    wrap command line
`,
	Action: doWrap,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "monitored `check-name` which must be unique on a host"},
		cli.BoolFlag{Name: "verbose, v", Usage: "verbose output"},
		cli.StringFlag{Name: "memo, m", Value: "", Usage: "`memo` of the job"},
		cli.StringFlag{Name: "H, host", Value: "", Usage: "`hostID`"},
		cli.BoolFlag{Name: "warning, w", Usage: "alerts as warning"},
	},
}

func doWrap(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	conf, err := config.LoadConfig(confFile)
	if err != nil {
		return err
	}
	apibase := c.GlobalString("apibase")
	if apibase == "" {
		apibase = conf.Apibase
	}

	apikey := conf.Apikey
	if apikey == "" {
		apikey = os.Getenv("MACKEREL_APIKEY")
	}
	if apikey == "" {
		logger.Log("error", "[mkr wrap] failed to detect Mackerel APIKey. Try to specify in mackerel-agent.conf or export MACKEREL_APIKEY='<Your apikey>'")
	}
	hostID, _ := conf.LoadHostID()
	if c.String("host") != "" {
		hostID = c.String("host")
	}
	if hostID == "" {
		logger.Log("error", "[mkr wrap] failed to load hostID. Try to specify -host option explicitly")
	}
	// Since command execution has the highest priority, even when apikey or
	// hostID is empty, we don't return errors and only output the log here.

	cmd := c.Args()
	if len(cmd) > 0 && cmd[0] == "--" {
		cmd = cmd[1:]
	}
	if len(cmd) < 1 {
		return fmt.Errorf("no commands specified")
	}

	return (&wrap{
		apibase: apibase,
		name:    c.String("name"),
		verbose: c.Bool("verbose"),
		memo:    c.String("memo"),
		warning: c.Bool("warning"),
		hostID:  hostID,
		apikey:  apikey,
		cmd:     cmd,
	}).run()
}
