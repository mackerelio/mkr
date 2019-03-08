package hosts

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/mackerelio/mkr/mackerelclient"
)

// Command is definition of mkr hosts subcommand
var Command = cli.Command{
	Name:      "hosts",
	Usage:     "List hosts",
	ArgsUsage: "[--verbose | -v] [--name | -n <name>] [--service | -s <service>] [[--role | -r <role>]...] [[--status | --st <status>]...]",
	Description: `
    List the information of the hosts refined by host name, service name, role name and/or status.
    Requests "GET /api/v0/hosts.json". See https://mackerel.io/api-docs/entry/hosts#list .
`,
	Action: doHosts,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "List hosts only matched with <name>"},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "List hosts only belonging to <service>"},
		cli.StringSliceFlag{
			Name:  "role, r",
			Value: &cli.StringSlice{},
			Usage: "List hosts only belonging to <role>. Multiple choices are allowed. Required --service",
		},
		cli.StringSliceFlag{
			Name:  "status, st",
			Value: &cli.StringSlice{},
			Usage: "List hosts only matched <status>. Multiple choices are allowed.",
		},
		cli.StringFlag{Name: "format, f", Value: "", Usage: "Output format template"},
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
	},
}

func doHosts(c *cli.Context) error {
	cli, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&hostApp{
		cli: cli,

		verbose: c.Bool("verbose"),

		name:     c.String("name"),
		service:  c.String("service"),
		roles:    c.StringSlice("role"),
		statuses: c.StringSlice("status"),

		format: c.String("format"),

		outStream: os.Stdout,
	}).run()
}
