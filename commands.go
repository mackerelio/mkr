package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Songmu/prompter"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/plugin"
	"gopkg.in/urfave/cli.v1"
)

func init() {
	// Requirements:
	// - .Description: First and last line is blank.
	// - .ArgsUsage: ArgsUsage includes flag usages (e.g. [-v|verbose] <hostId>).
	//   All cli.Command should have ArgsUsage field.
	cli.CommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Description}}

DESCRIPTION:{{.Description}}{{end}}{{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
}

// Commands cli.Command object list
var Commands = []cli.Command{
	commandStatus,
	commandHosts,
	commandCreate,
	commandUpdate,
	commandThrow,
	commandMetrics,
	commandFetch,
	commandRetire,
	commandServices,
	commandMonitors,
	commandAlerts,
	commandDashboards,
	commandAnnotations,
	commandOrg,
	plugin.CommandPlugin,
}

var commandStatus = cli.Command{
	Name:      "status",
	Usage:     "Show the host",
	ArgsUsage: "[--verbose | -v] <hostId>",
	Description: `
    Show the information of the host identified with <hostId>.
    Requests "GET /api/v0/hosts/<hostId>". See https://mackerel.io/api-docs/entry/hosts#get .
`,
	Action: doStatus,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
	},
}

var commandHosts = cli.Command{
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

var commandCreate = cli.Command{
	Name:      "create",
	Usage:     "Create a new host",
	ArgsUsage: "[--status | -st <status>] [--roleFullname | -R <service:role>] [--customIdentifier <customIdentifier>] <hostName>",
	Description: `
    Create a new host with status, roleFullname and/or customIdentifier.
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
	},
}

var commandUpdate = cli.Command{
	Name:      "update",
	Usage:     "Update the host",
	ArgsUsage: "[--name | -n <name>] [--displayName <displayName>] [--status | -st <status>] [--roleFullname | -R <service:role>] [--overwriteRoles | -o] [<hostIds...>]",
	Description: `
    Update the host identified with <hostId>.
    Requests "PUT /api/v0/hosts/<hostId>". See https://mackerel.io/api-docs/entry/hosts#update-information .
`,
	Action: doUpdate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "Update hostname."},
		cli.StringFlag{Name: "displayName", Value: "", Usage: "Update displayName."},
		cli.StringFlag{Name: "status, st", Value: "", Usage: "Update status."},
		cli.StringSliceFlag{
			Name:  "roleFullname, R",
			Value: &cli.StringSlice{},
			Usage: "Update rolefullname.",
		},
		cli.BoolFlag{Name: "overwriteRoles, o", Usage: "Overwrite roles instead of adding specified roles."},
	},
}

var commandMetrics = cli.Command{
	Name:      "metrics",
	Usage:     "Fetch metric values",
	ArgsUsage: "[--host | -H <hostId>] [--service | -s <service>] [--name | -n <metricName>] --from int --to int",
	Description: `
    Fetch metric values of 'host metric' or 'service metric'.
    Requests "/api/v0/hosts/<hostId>/metrics" or "/api/v0/services/<serviceName>/tsdb".
		See https://mackerel.io/api-docs/entry/host-metrics#get, https://mackerel.io/ja/api-docs/entry/service-metrics#get.
`,
	Action: doMetrics,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host, H", Value: "", Usage: "Fetch host metric values of <hostID>."},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Fetch service metric values of <service>."},
		cli.StringFlag{Name: "name, n", Value: "", Usage: "The name of the metric for which you want to obtain the metric."},
		cli.Int64Flag{Name: "from", Usage: "The first of the period for which you want to obtain the metric. (epoch seconds)"},
		cli.Int64Flag{Name: "to", Usage: "The end of the period for which you want to obtain the metric. (epoch seconds)"},
	},
}

var commandFetch = cli.Command{
	Name:      "fetch",
	Usage:     "Fetch latest metric values",
	ArgsUsage: "[--name | -n <metricName>] hostIds...",
	Description: `
    Fetch latest metric values about the hosts.
    Requests "GET /api/v0/tsdb/latest". See https://mackerel.io/api-docs/entry/host-metrics#get-latest .
`,
	Action: doFetch,
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "name, n",
			Value: &cli.StringSlice{},
			Usage: "Fetch metric values identified with <name>. Required. Multiple choices are allowed. ",
		},
	},
}

var commandRetire = cli.Command{
	Name:      "retire",
	Usage:     "Retire hosts",
	ArgsUsage: "[--force] hostIds...",
	Description: `
    Retire host identified by <hostId>. Be careful because this is an irreversible operation.
    Requests POST /api/v0/hosts/<hostId>/retire parallelly. See https://mackerel.io/api-docs/entry/hosts#retire .
`,
	Action: doRetire,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "force", Usage: "Force retirement without confirmation."},
	},
}

var commandServices = cli.Command{
	Name:      "services",
	Usage:     "List services",
	ArgsUsage: "",
	Description: `
    List the information of the services.
    Requests "GET /api/v0/services". See https://mackerel.io/api-docs/entry/services#list.
`,
	Action: doServices,
	Flags:  []cli.Flag{},
}

func newMackerelFromContext(c *cli.Context) *mkr.Client {
	confFile := c.GlobalString("conf")
	apiBase := c.GlobalString("apibase")
	apiKey := LoadApikeyFromEnvOrConfig(confFile)
	if apiKey == "" {
		logger.Log("error", `
    MACKEREL_APIKEY environment variable is not set. (Try "export MACKEREL_APIKEY='<Your apikey>'")
`)
		os.Exit(1)
	}

	if apiBase == "" {
		apiBase = LoadApibaseFromConfigWithFallback(confFile)
	}

	mackerel, err := mkr.NewClientWithOptions(apiKey, apiBase, os.Getenv("DEBUG") != "")
	logger.DieIf(err)

	return mackerel
}

func doStatus(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	argHostID := c.Args().Get(0)
	isVerbose := c.Bool("verbose")

	if argHostID == "" {
		if argHostID = LoadHostIDFromConfig(confFile); argHostID == "" {
			cli.ShowCommandHelp(c, "status")
			os.Exit(1)
		}
	}

	host, err := newMackerelFromContext(c).FindHost(argHostID)
	logger.DieIf(err)

	if isVerbose {
		PrettyPrintJSON(host)
	} else {
		format := &HostFormat{
			ID:            host.ID,
			Name:          host.Name,
			DisplayName:   host.DisplayName,
			Status:        host.Status,
			RoleFullnames: host.GetRoleFullnames(),
			IsRetired:     host.IsRetired,
			CreatedAt:     formatISO8601Extended(host.DateFromCreatedAt()),
			IPAddresses:   host.IPAddresses(),
		}

		PrettyPrintJSON(format)
	}
	return nil
}

func doHosts(c *cli.Context) error {
	isVerbose := c.Bool("verbose")

	hosts, err := newMackerelFromContext(c).FindHosts(&mkr.FindHostsParam{
		Name:     c.String("name"),
		Service:  c.String("service"),
		Roles:    c.StringSlice("role"),
		Statuses: c.StringSlice("status"),
	})
	logger.DieIf(err)

	format := c.String("format")
	if format != "" {
		t := template.Must(template.New("format").Parse(format))
		err := t.Execute(os.Stdout, hosts)
		logger.DieIf(err)
	} else if isVerbose {
		PrettyPrintJSON(hosts)
	} else {
		var hostsFormat []*HostFormat
		for _, host := range hosts {
			format := &HostFormat{
				ID:            host.ID,
				Name:          host.Name,
				DisplayName:   host.DisplayName,
				Status:        host.Status,
				RoleFullnames: host.GetRoleFullnames(),
				IsRetired:     host.IsRetired,
				CreatedAt:     formatISO8601Extended(host.DateFromCreatedAt()),
				IPAddresses:   host.IPAddresses(),
			}
			hostsFormat = append(hostsFormat, format)
		}

		PrettyPrintJSON(hostsFormat)
	}
	return nil
}

func doCreate(c *cli.Context) error {
	argHostName := c.Args().Get(0)
	optRoleFullnames := c.StringSlice("roleFullname")
	optStatus := c.String("status")
	optCustomIdentifier := c.String("customIdentifier")

	if argHostName == "" {
		cli.ShowCommandHelp(c, "create")
		os.Exit(1)
	}

	client := newMackerelFromContext(c)

	hostID, err := client.CreateHost(&mkr.CreateHostParam{
		Name:             argHostName,
		RoleFullnames:    optRoleFullnames,
		CustomIdentifier: optCustomIdentifier,
	})
	logger.DieIf(err)

	logger.Log("created", hostID)

	if optStatus != "" {
		err := client.UpdateHostStatus(hostID, optStatus)
		logger.DieIf(err)
		logger.Log("updated", fmt.Sprintf("%s %s", hostID, optStatus))
	}
	return nil
}

func doUpdate(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	argHostIDs := c.Args()
	optName := c.String("name")
	optDisplayName := c.String("displayName")
	optStatus := c.String("status")
	optRoleFullnames := c.StringSlice("roleFullname")
	overwriteRoles := c.Bool("overwriteRoles")

	if len(argHostIDs) < 1 {
		argHostIDs = make([]string, 1)
		if argHostIDs[0] = LoadHostIDFromConfig(confFile); argHostIDs[0] == "" {
			cli.ShowCommandHelp(c, "update")
			os.Exit(1)
		}
	}

	needUpdateHostStatus := optStatus != ""
	needUpdateRolesInHostUpdate := !overwriteRoles && len(optRoleFullnames) > 0
	needUpdateHost := (optName != "" || optDisplayName != "" || overwriteRoles || needUpdateRolesInHostUpdate)

	if !needUpdateHostStatus && !needUpdateHost {
		logger.Log("update", "at least one argumet is required.")
		cli.ShowCommandHelp(c, "update")
		os.Exit(1)
	}

	client := newMackerelFromContext(c)

	for _, hostID := range argHostIDs {
		if needUpdateHostStatus {
			err := client.UpdateHostStatus(hostID, optStatus)
			logger.DieIf(err)
		}

		if overwriteRoles {
			err := client.UpdateHostRoleFullnames(hostID, optRoleFullnames)
			logger.DieIf(err)
		}

		if needUpdateHost {
			host, err := client.FindHost(hostID)
			logger.DieIf(err)
			meta := host.Meta
			name := ""
			if optName == "" {
				name = host.Name
			} else {
				name = optName
			}
			displayname := ""
			if optDisplayName == "" {
				displayname = host.DisplayName
			} else {
				displayname = optDisplayName
			}
			param := &mkr.UpdateHostParam{
				Name:        name,
				DisplayName: displayname,
				Meta:        meta,
			}
			if needUpdateRolesInHostUpdate {
				param.RoleFullnames = optRoleFullnames
			}
			_, err = client.UpdateHost(hostID, param)
			logger.DieIf(err)
		}

		logger.Log("updated", hostID)
	}
	return nil
}

func split(ids []string, count int) [][]string {
	xs := make([][]string, 0, (len(ids)+count-1)/count)
	for i, name := range ids {
		if i%count == 0 {
			xs = append(xs, []string{})
		}
		xs[len(xs)-1] = append(xs[len(xs)-1], name)
	}
	return xs
}

func doMetrics(c *cli.Context) error {
	optHostID := c.String("host")
	optService := c.String("service")
	optMetricName := c.String("name")

	from := c.Int64("from")
	to := c.Int64("to")
	if to == 0 {
		to = time.Now().Unix()
	}

	client := newMackerelFromContext(c)

	if optHostID != "" {
		metricValue, err := client.FetchHostMetricValues(optHostID, optMetricName, from, to)
		logger.DieIf(err)

		PrettyPrintJSON(metricValue)
	} else if optService != "" {
		metricValue, err := client.FetchServiceMetricValues(optService, optMetricName, from, to)
		logger.DieIf(err)

		PrettyPrintJSON(metricValue)
	} else {
		cli.ShowCommandHelp(c, "metrics")
		os.Exit(1)
	}
	return nil
}

func doFetch(c *cli.Context) error {
	argHostIDs := c.Args()
	optMetricNames := c.StringSlice("name")

	if len(argHostIDs) < 1 || len(optMetricNames) < 1 {
		cli.ShowCommandHelp(c, "fetch")
		os.Exit(1)
	}

	allMetricValues := make(mkr.LatestMetricValues)
	// Fetches 100 hosts per one request (to avoid URL maximum length).
	for _, hostIds := range split(argHostIDs, 100) {
		metricValues, err := newMackerelFromContext(c).FetchLatestMetricValues(hostIds, optMetricNames)
		logger.DieIf(err)
		for key := range metricValues {
			allMetricValues[key] = metricValues[key]
		}
	}

	PrettyPrintJSON(allMetricValues)
	return nil
}

func doRetire(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	force := c.Bool("force")
	argHostIDs := c.Args()

	if len(argHostIDs) < 1 {
		argHostIDs = make([]string, 1)
		if argHostIDs[0] = LoadHostIDFromConfig(confFile); argHostIDs[0] == "" {
			cli.ShowCommandHelp(c, "retire")
			os.Exit(1)
		}
	}

	if !force && !prompter.YN("Retire following hosts.\n  "+strings.Join(argHostIDs, "\n  ")+"\nAre you sure?", true) {
		logger.Log("", "retirement is canceled.")
		return nil
	}

	client := newMackerelFromContext(c)

	for _, hostID := range argHostIDs {
		err := client.RetireHost(hostID)
		logger.DieIf(err)

		logger.Log("retired", hostID)
	}
	return nil
}

func doServices(c *cli.Context) error {
	services, err := newMackerelFromContext(c).FindServices()
	logger.DieIf(err)
	PrettyPrintJSON(services)
	return nil
}
