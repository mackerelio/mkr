package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/Songmu/prompter"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
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
	commandFetch,
	commandRetire,
	commandServices,
	commandMonitors,
	commandAlerts,
	commandDashboards,
	commandAnnotations,
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
	ArgsUsage: "[--status | -st <status>] [--roleFullname | -R <service:role>] <hostName>",
	Description: `
    Create a new host with status and/or roleFullname.
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

var commandThrow = cli.Command{
	Name:      "throw",
	Usage:     "Post metric values",
	ArgsUsage: "[--host | -h <hostId>] [--service | -s <service>] stdin",
	Description: `
    Post metric values to 'host metric' or 'service metric'.
    Output format of metric values are compatible with that of a Sensu plugin.
    Requests "POST /api/v0/tsdb". See https://mackerel.io/api-docs/entry/host-metrics#post .
`,
	Action: doThrow,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host, H", Value: "", Usage: "Post host metric values to <hostID>."},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Post service metric values to <service>."},
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
			CreatedAt:     host.DateStringFromCreatedAt(),
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
				CreatedAt:     host.DateStringFromCreatedAt(),
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

	if argHostName == "" {
		cli.ShowCommandHelp(c, "create")
		os.Exit(1)
	}

	client := newMackerelFromContext(c)

	hostID, err := client.CreateHost(&mkr.CreateHostParam{
		Name:          argHostName,
		RoleFullnames: optRoleFullnames,
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

func doThrow(c *cli.Context) error {
	optHostID := c.String("host")
	optService := c.String("service")

	var metricValues []*(mkr.MetricValue)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// name, value, timestamp
		// ex.) tcp.CLOSING 0 1397031808
		items := strings.Fields(line)
		if len(items) != 3 {
			continue
		}
		value, err := strconv.ParseFloat(items[1], 64)
		if err != nil {
			logger.Log("warning", fmt.Sprintf("Failed to parse values: %s", err))
			continue
		}
		time, err := strconv.ParseInt(items[2], 10, 64)
		if err != nil {
			logger.Log("warning", fmt.Sprintf("Failed to parse values: %s", err))
			continue
		}

		name := items[0]
		if optHostID != "" && !strings.HasPrefix(name, "custom.") {
			name = "custom." + name
		}

		metricValue := &mkr.MetricValue{
			Name:  name,
			Value: value,
			Time:  time,
		}

		metricValues = append(metricValues, metricValue)
	}
	logger.ErrorIf(scanner.Err())

	client := newMackerelFromContext(c)

	if optHostID != "" {
		err := client.PostHostMetricValuesByHostID(optHostID, metricValues)
		logger.DieIf(err)

		for _, metric := range metricValues {
			logger.Log("thrown", fmt.Sprintf("%s '%s\t%f\t%d'", optHostID, metric.Name, metric.Value, metric.Time))
		}
	} else if optService != "" {
		err := client.PostServiceMetricValues(optService, metricValues)
		logger.DieIf(err)

		for _, metric := range metricValues {
			logger.Log("thrown", fmt.Sprintf("%s '%s\t%f\t%d'", optService, metric.Name, metric.Value, metric.Time))
		}
	} else {
		cli.ShowCommandHelp(c, "throw")
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

	metricValues, err := newMackerelFromContext(c).FetchLatestMetricValues(argHostIDs, optMetricNames)
	logger.DieIf(err)

	PrettyPrintJSON(metricValues)
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
