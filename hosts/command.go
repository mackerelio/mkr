package hosts

import (
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v2"
)

// CommandCreate is definition of mkr create subcommand
var CommandCreate = cli.Command{
	Name:      "create",
	Usage:     "Create a new host",
	ArgsUsage: "[--status | -st <status>] [--roleFullname | -R <service:role>] [--customIdentifier <customIdentifier>] [--memo <memo>] <hostName>",
	Description: `
    Create a new host with status, roleFullname, customIdentifier and/or memo.
    Requests "POST /api/v0/hosts". See https://mackerel.io/api-docs/entry/hosts#create .
`,
	Action: doCreate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "status, st", Value: "", Usage: "Host status ('working', 'standby', 'maintenance')"},
		cli.StringSliceFlag{
			Name:  "roleFullname, R",
			Value: &cli.StringSlice{},
			Usage: "Multiple choices are allowed. ex. My-Service:proxy, My-Service:db-master",
		},
		cli.StringFlag{Name: "customIdentifier", Value: "", Usage: "CustomIdentifier for the Host"},
		cli.StringFlag{Name: "memo", Value: "", Usage: "memo for the Host"},
	},
}

func doCreate(c *cli.Context) error {
	argHostName := c.Args().Get(0)
	if argHostName == "" {
		cli.ShowCommandHelpAndExit(c, "create", 1)
	}

	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&hostApp{
		client:    client,
		logger:    logger.New(),
		outStream: os.Stdout,
	}).createHost(createHostParam{
		name:             argHostName,
		roleFullnames:    c.StringSlice("roleFullname"),
		status:           c.String("status"),
		customIdentifier: c.String("customIdentifier"),
		memo:             c.String("memo"),
	})
}

// CommandHosts is definition of mkr hosts subcommand
var CommandHosts = cli.Command{
	Name:      "hosts",
	Usage:     "List hosts",
	ArgsUsage: "[--verbose | -v] [--name | -n <name>] [--service | -s <service>] [[--role | -r <role>]...] [[--status | --st <status>]...] [--jq <formula>]",
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
		jq.CommandLineFlag,
	},
}

func doHosts(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&hostApp{
		client:    client,
		logger:    logger.New(),
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).findHosts(findHostsParam{
		verbose: c.Bool("verbose"),

		name:     c.String("name"),
		service:  c.String("service"),
		roles:    c.StringSlice("role"),
		statuses: c.StringSlice("status"),

		format: c.String("format"),
	})
}
