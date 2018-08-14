package wrap

import (
	"gopkg.in/urfave/cli.v1"
)

// CommandPlugin is definition of mkr plugin
var Command = cli.Command{
	Name:  "wrap",
	Usage: "wrap command status",
	ArgsUsage: "[--verbose | -v] [--name | -n <name>] [--memo | -m <memo>] -- /path/to/cmd",
	Description: `
	wrap command line
`,
	Action: doWrap,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "monitor <name>"},
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
		cli.StringFlag{Name: "memo, m", Value: "", Usage: "monitor <memo>"},
	},
}

func doWrap(c *cli.Context) error {
	return nil
}
