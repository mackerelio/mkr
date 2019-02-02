package hosts

import (
	"os"
	"text/template"

	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	cli "gopkg.in/urfave/cli.v1"
)

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
	isVerbose := c.Bool("verbose")

	hosts, err := mackerelclient.NewFromContext(c).FindHosts(&mkr.FindHostsParam{
		Name:     c.String("name"),
		Service:  c.String("service"),
		Roles:    c.StringSlice("role"),
		Statuses: c.StringSlice("status"),
	})
	logger.DieIf(err)

	fmtStr := c.String("format")
	if fmtStr != "" {
		t := template.Must(template.New("format").Parse(fmtStr))
		err := t.Execute(os.Stdout, hosts)
		logger.DieIf(err)
	} else if isVerbose {
		format.PrettyPrintJSON(hosts)
	} else {
		var hostsFormat []*format.Host
		for _, host := range hosts {
			hostsFormat = append(hostsFormat, &format.Host{
				ID:            host.ID,
				Name:          host.Name,
				DisplayName:   host.DisplayName,
				Status:        host.Status,
				RoleFullnames: host.GetRoleFullnames(),
				IsRetired:     host.IsRetired,
				CreatedAt:     format.ISO8601Extended(host.DateFromCreatedAt()),
				IPAddresses:   host.IPAddresses(),
			})
		}

		format.PrettyPrintJSON(hostsFormat)
	}
	return nil
}
