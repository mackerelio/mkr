package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
)

// cli.Commad object list
var Commands = []cli.Command{
	commandStatus,
	commandHosts,
	commandCreate,
	commandUpdate,
	commandThrow,
	commandFetch,
	commandRetire,
}

var commandStatus = cli.Command{
	Name:  "status",
	Usage: "Show the host",
	Description: `
    Show the information of the host identified with <hostId>.
    Request "GET /api/v0/hosts/<hostId>". See http://help-ja.mackerel.io/entry/spec/api/v0#host-get.
`,
	Action: doStatus,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
	},
}

var commandHosts = cli.Command{
	Name:  "hosts",
	Usage: "List hosts",
	Description: `
    List the information of the hosts refined by host name, service name, role name and/or status.
    Request "GET /api/v0/hosts.json". See http://help-ja.mackerel.io/entry/spec/api/v0#hosts-list.
`,
	Action: doHosts,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "List hosts only matched with <name>"},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "List hosts only belongs to <service>"},
		cli.StringSliceFlag{
			Name:  "role, r",
			Value: &cli.StringSlice{},
			Usage: "List hosts only belongs to <role>. Multiple choice allow. Required --service",
		},
		cli.StringSliceFlag{
			Name:  "status, st",
			Value: &cli.StringSlice{},
			Usage: "List hosts only matched <status>. Multiple choice allow.",
		},
		cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
	},
}

var commandCreate = cli.Command{
	Name:  "create",
	Usage: "Create a new host",
	Description: `
    Create a new host with staus and/or roleFullname.
    Request "POST /api/v0/hosts". See http://help-ja.mackerel.io/entry/spec/api/v0#host-create.
`,
	Action: doCreate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "status, st", Value: "", Usage: "Host status ('working', 'standby', 'meintenance')"},
		cli.StringSliceFlag{
			Name:  "roleFullname, R",
			Value: &cli.StringSlice{},
			Usage: "Multiple choice allow. ex. My-Service:proxy, My-Service:db-master",
		},
	},
}

var commandUpdate = cli.Command{
	Name:  "update",
	Usage: "Update the host",
	Description: `
    Update the hosts identified with <hostIds>.
    Request "PUT /api/v0/hosts/<hostId>". See http://help-ja.mackerel.io/entry/spec/api/v0#host-update.
`,
	Action: doUpdate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "Update hostname."},
		cli.StringFlag{Name: "status, st", Value: "", Usage: "Update status."},
		cli.StringSliceFlag{
			Name:  "roleFullname, R",
			Value: &cli.StringSlice{},
			Usage: "Update rolefullname.",
		},
	},
}

var commandThrow = cli.Command{
	Name:  "throw",
	Usage: "Post metric values",
	Description: `
    Post metric values to 'host metric' or 'service metric'.
    Output format of metric value is compatibled with that of Sensu plugin.
    Request "POST /api/v0/tsdb". See http://help-ja.mackerel.io/entry/spec/api/v0#metric-value-post.
`,
	Action: doThrow,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host, H", Value: "", Usage: "Post host metric values to <hostID>."},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Post service metric values to <service>."},
	},
}

var commandFetch = cli.Command{
	Name:  "fetch",
	Usage: "Fetch latest metric values",
	Description: `
    Fetch latest metric values about the hosts.
    Request "GET /api/v0/tsdb/latest". See http://help-ja.mackerel.io/entry/spec/api/v0#tsdb-latest.
`,
	Action: doFetch,
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "name, n",
			Value: &cli.StringSlice{},
			Usage: "Fetch metric values identified with <name>. Required. Multiple choise allow. ",
		},
	},
}

var commandRetire = cli.Command{
	Name:  "retire",
	Usage: "Retire hosts",
	Description: `
    Retire host identified by hostIds. Be careful because this is a irreversible operation.
    Request POST /api/v0/hosts/<hostId>/retire parallelly. See http://help-ja.mackerel.io/entry/spec/api/v0#host-retire.
`,
	Action: doRetire,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newMackerel() *mkr.Client {
	apiKey := LoadApikeyFromEnvOrConfig()
	if apiKey == "" {
		logger.Log("error", `
    Not set MACKEREL_APIKEY environment variable. (Try "export MACKEREL_APIKEY='<Your apikey>'")
`)
		os.Exit(1)
	}

	if os.Getenv("DEBUG") != "" {
		mackerel, err := mkr.NewClientWithOptions(apiKey, "https://mackerel.io/api/v0", true)
		logger.DieIf(err)

		return mackerel
	}

	return mkr.NewClient(apiKey)
}

type commandDoc struct {
	Parent    string
	Arguments string
}

var commandDocs = map[string]commandDoc{
	"status": {"", "[-v|verbose] <hostId>"},
	"hosts":  {"", "[--verbose | -v] [--name | -n <name>] [--service | -s <service>] [[--role | -r <role>]...] [[--status | --st <status>]...]"},
	"create": {"", "[--status | -st <status>] [--roleFullname | -R <service:role>] <hostName>"},
	"update": {"", "[--name | -n <name>] [--status | -st <status>] [--roleFullname | -R <service:role>] <hostIds...> ]"},
	"throw":  {"", "[--host | -h <hostId>] [--service | -s <service>] stdin"},
	"fetch":  {"", "[--name | -n <metricName>] hostIds..."},
	"retire": {"", "hostIds..."},
}

// Makes template conditionals to generate per-command documents.
func mkCommandsTemplate(genTemplate func(commandDoc) string) string {
	template := "{{if false}}"
	for _, command := range append(Commands) {
		template = template + fmt.Sprintf("{{else if (eq .Name %q)}}%s", command.Name, genTemplate(commandDocs[command.Name]))
	}
	return template + "{{end}}"
}

func init() {
	argsTemplate := mkCommandsTemplate(func(doc commandDoc) string { return doc.Arguments })
	parentTemplate := mkCommandsTemplate(func(doc commandDoc) string { return string(strings.TrimLeft(doc.Parent+" ", " ")) })

	cli.CommandHelpTemplate = `NAME:
    {{.Name}} - {{.Usage}}

USAGE:
    mkr ` + parentTemplate + `{{.Name}} ` + argsTemplate + `
{{if (len .Description)}}
DESCRIPTION: {{.Description}}
{{end}}{{if (len .Flags)}}
OPTIONS:
    {{range .Flags}}{{.}}
    {{end}}
{{end}}`
}

func doStatus(c *cli.Context) {
	argHostID := c.Args().Get(0)
	isVerbose := c.Bool("verbose")

	if argHostID == "" {
		if argHostID = LoadHostIDFromConfig(); argHostID == "" {
			cli.ShowCommandHelp(c, "status")
			os.Exit(1)
		}
	}

	host, err := newMackerel().FindHost(argHostID)
	logger.DieIf(err)

	if isVerbose {
		PrettyPrintJSON(host)
	} else {
		format := &HostFormat{
			ID:            host.ID,
			Name:          host.Name,
			Status:        host.Status,
			RoleFullnames: host.GetRoleFullnames(),
			IsRetired:     host.IsRetired,
			CreatedAt:     host.DateStringFromCreatedAt(),
			IPAddresses:   host.IPAddresses(),
		}

		PrettyPrintJSON(format)
	}
}

func doHosts(c *cli.Context) {
	isVerbose := c.Bool("verbose")

	hosts, err := newMackerel().FindHosts(&mkr.FindHostsParam{
		Name:     c.String("name"),
		Service:  c.String("service"),
		Roles:    c.StringSlice("role"),
		Statuses: c.StringSlice("status"),
	})
	logger.DieIf(err)

	if isVerbose {
		PrettyPrintJSON(hosts)
	} else {
		var hostsFormat []*HostFormat
		for _, host := range hosts {
			format := &HostFormat{
				ID:            host.ID,
				Name:          host.Name,
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
}

func doCreate(c *cli.Context) {
	argHostName := c.Args().Get(0)
	optRoleFullnames := c.StringSlice("roleFullname")
	optStatus := c.String("status")

	if argHostName == "" {
		cli.ShowCommandHelp(c, "create")
		os.Exit(1)
	}

	client := newMackerel()

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
}

func doUpdate(c *cli.Context) {
	argHostIDs := c.Args()
	optName := c.String("name")
	optStatus := c.String("status")
	optRoleFullnames := c.StringSlice("roleFullname")

	if len(argHostIDs) < 1 {
		argHostIDs = make([]string, 1)
		if argHostIDs[0] = LoadHostIDFromConfig(); argHostIDs[0] == "" {
			cli.ShowCommandHelp(c, "update")
			os.Exit(1)
		}
	}

	needUpdateHostStatus := optStatus != ""
	needUpdateHost := (optName != "" || len(optRoleFullnames) > 0)

	if !needUpdateHostStatus && !needUpdateHost {
		cli.ShowCommandHelp(c, "update")
		os.Exit(1)
	}

	client := newMackerel()

	var wg sync.WaitGroup
	for _, hostID := range argHostIDs {
		wg.Add(1)
		go func(hostID string) {
			defer wg.Done()

			if needUpdateHostStatus {
				err := client.UpdateHostStatus(hostID, optStatus)
				logger.DieIf(err)
			}

			if needUpdateHost {
				_, err := client.UpdateHost(hostID, &mkr.UpdateHostParam{
					Name:          optName,
					RoleFullnames: optRoleFullnames,
				})
				logger.DieIf(err)
			}

			logger.Log("updated", hostID)
		}(hostID)
	}

	wg.Wait()
}

func doThrow(c *cli.Context) {
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

		metricValue := &mkr.MetricValue{
			Name:  items[0],
			Value: value,
			Time:  time,
		}

		metricValues = append(metricValues, metricValue)
	}
	logger.ErrorIf(scanner.Err())

	client := newMackerel()

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
}

func doFetch(c *cli.Context) {
	argHostIDs := c.Args()
	optMetricNames := c.StringSlice("name")

	if len(argHostIDs) < 1 || len(optMetricNames) < 1 {
		cli.ShowCommandHelp(c, "fetch")
		os.Exit(1)
	}

	metricValues, err := newMackerel().FetchLatestMetricValues(argHostIDs, optMetricNames)
	logger.DieIf(err)

	PrettyPrintJSON(metricValues)
}

func doRetire(c *cli.Context) {
	argHostIDs := c.Args()

	if len(argHostIDs) < 1 {
		argHostIDs = make([]string, 1)
		if argHostIDs[0] = LoadHostIDFromConfig(); argHostIDs[0] == "" {
			cli.ShowCommandHelp(c, "retire")
			os.Exit(1)
		}
	}

	client := newMackerel()

	var wg sync.WaitGroup
	for _, hostID := range argHostIDs {
		wg.Add(1)
		go func(hostID string) {
			defer wg.Done()

			err := client.RetireHost(hostID)
			logger.DieIf(err)

			logger.Log("retired", hostID)
		}(hostID)
	}

	wg.Wait()
}
