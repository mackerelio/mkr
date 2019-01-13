package wrap

import (
	"fmt"
	"os"

	"github.com/mackerelio/mackerel-agent/config"
	cli "gopkg.in/urfave/cli.v1"
)

// CommandPlugin is definition of mkr plugin
var Command = cli.Command{
	Name:      "wrap",
	Usage:     "wrap command status",
	ArgsUsage: "[--verbose | -v] [--name | -n <name>] [--memo | -m <memo>] -- /path/to/batch",
	Description: `
    wrap command line
`,
	Action: doWrap,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "monitor <name>"},
		cli.BoolFlag{Name: "verbose, v", Usage: "verbose output mode"},
		cli.StringFlag{Name: "memo, m", Value: "", Usage: "monitor <memo>"},
		cli.StringFlag{Name: "H, host", Value: "", Usage: "<hostId>"},
		cli.BoolFlag{Name: "warning, w", Usage: "alert as warning"},
	},
}

func doWrap(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	conf, err := config.LoadConfig(confFile)
	if err != nil {
		return err
	}
	apibase := c.GlobalString("apibase")
	apikey := conf.Apikey
	if apikey == "" {
		apikey = os.Getenv("MACKEREL_APIKEY")
	}
	if apikey == "" {
		return fmt.Errorf(`MACKEREL_APIKEY environment variable is not set. (Try "export MACKEREL_APIKEY='<Your apikey>'`)
	}
	hostID, _ := conf.LoadHostID()
	if c.String("host") != "" {
		hostID = c.String("host")
	}
	if hostID == "" {
		return fmt.Errorf("failed to load hostID. (Try to specify -host option explicitly)")
	}
	cmd := c.Args()
	if len(cmd) > 0 && cmd[0] == "--" {
		cmd = cmd[1:]
	}
	if len(cmd) < 1 {
		return fmt.Errorf("no command specified")
	}

	return (&app{
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

type app struct {
	apibase string
	name    string
	verbose bool
	memo    string
	warning bool
	hostID  string
	apikey  string
	cmd     []string
}

func (ap *app) run() error {
	return nil
}
