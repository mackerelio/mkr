package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mackerelio/gomkr/utils"
	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
	mkr "github.com/mackerelio/mackerel-client-go"
)

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
	Usage: "Show a host",
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
		cli.StringFlag{Name: "name, n", Value: "", Usage: "Show hosts only matched with <name>"},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Show hosts only belongs to <service>"},
		cli.StringSliceFlag{
			Name:  "role, r",
			Value: &cli.StringSlice{},
			Usage: "Show hosts only belongs to <role>. Multiple choice allow. Required --service",
		},
		cli.StringSliceFlag{
			Name:  "status, st",
			Value: &cli.StringSlice{},
			Usage: "Show hosts only matched <status>. Multiple choice allow.",
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
    Update the host identified with <hostId>.
    Request "PUT /api/v0/hosts/<hostId>". See http://help-ja.mackerel.io/entry/spec/api/v0#host-update.
`,
	Action: doUpdate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "Update <hostId> hostname to <name>."},
		cli.StringFlag{Name: "status, st", Value: "", Usage: "Update <hostId> status to <status>."},
		cli.StringSliceFlag{
			Name:  "roleFullname, R",
			Value: &cli.StringSlice{},
			Usage: "Update <hostId> rolefullname to <roleFullname>.",
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
		cli.StringFlag{Name: "host, H", Value: "", Usage: "Post host metric values to <hostId>."},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Post service metric values to <service>."},
	},
}

var commandFetch = cli.Command{
	Name:  "fetch",
	Usage: "Fetch latest metric values",
	Description: `
    Fetch latest metric values about <hostId>... hosts.
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
	Usage: "Retire host",
	Description: `
    Retire host identified by <hostId>. Be careful because this is a irreversible operation.
    Request POST /api/v0/hosts/<hostId>/retire. See http://help-ja.mackerel.io/entry/spec/api/v0#host-retire.
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
	apiKey := os.Getenv("MACKEREL_APIKEY")
	if apiKey == "" {
		utils.Log("error", `
    Not set MACKEREL_APIKEY environment variable. (Try "export MACKEREL_APIKEY='<Your apikey>'")
`)
		os.Exit(1)
	}

	if os.Getenv("DEBUG") != "" {
		mackerel, err := mkr.NewClientForTest(apiKey, "https://mackerel.io/api/v0", true)
		utils.DieIf(err)

		return mackerel
	} else {
		return mkr.NewClient(apiKey)
	}
}

type commandDoc struct {
	Parent    string
	Arguments string
}

var commandDocs = map[string]commandDoc{
	"status": {"", "[-v|verbose]"},
	"hosts":  {"", "[--verbose | -v] [--name | -n <name>] [--service | -s <service>] [[--role | -r <role>]...] [[--status | --st <status>]...]"},
	"create": {"", "[--status | -st <status>] [--roleFullname | -R <service:role>] <hostName>"},
	"update": {"", "[--name | -n <name>] [--status | -st <status>] [--roleFullname | -R <service:role>] <hostId>"},
	"throw":  {"", "[--host | -h <hostId>] [--service | -s <service>] stdin"},
	"fetch":  {"", "[--name | -n <metricName>] <hostId>..."},
	"retire": {"", "<hostId>"},
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
    gomkr ` + parentTemplate + `{{.Name}} ` + argsTemplate + `
{{if (len .Description)}}
DESCRIPTION: {{.Description}}
{{end}}{{if (len .Flags)}}
OPTIONS:
    {{range .Flags}}{{.}}
    {{end}}
{{end}}`
}

func LoadHostIdFromConfig() string {
	conf, err := config.LoadConfig(config.DefaultConfig.Conffile)
	if err != nil {
		return ""
	}
	hostId, err := command.LoadHostId(conf.Root)
	if err != nil {
		return ""
	}
	return hostId
}

func doStatus(c *cli.Context) {
	argHostId := c.Args().Get(0)
	isVerbose := c.Bool("verbose")

	if argHostId == "" {
		if argHostId = LoadHostIdFromConfig(); argHostId == "" {
			cli.ShowCommandHelp(c, "status")
			os.Exit(1)
		}
	}

	host, err := newMackerel().FindHost(argHostId)
	utils.DieIf(err)

	if isVerbose {
		PrettyPrintJson(host)
	} else {
		format := &HostFormat{
			Id:            host.Id,
			Name:          host.Name,
			Status:        host.Status,
			RoleFullnames: host.GetRoleFullnames(),
			IsRetired:     host.IsRetired,
			CreatedAt:     host.DateStringFromCreatedAt(),
			IpAddresses:   host.IpAddresses(),
		}

		PrettyPrintJson(format)
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
	utils.DieIf(err)

	if isVerbose {
		PrettyPrintJson(hosts)
	} else {
		var hostsFormat []*HostFormat
		for _, host := range hosts {
			format := &HostFormat{
				Id:            host.Id,
				Name:          host.Name,
				Status:        host.Status,
				RoleFullnames: host.GetRoleFullnames(),
				IsRetired:     host.IsRetired,
				CreatedAt:     host.DateStringFromCreatedAt(),
				IpAddresses:   host.IpAddresses(),
			}
			hostsFormat = append(hostsFormat, format)
		}

		PrettyPrintJson(hostsFormat)
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

	hostId, err := newMackerel().CreateHost(&mkr.CreateHostParam{
		Name:          argHostName,
		RoleFullnames: optRoleFullnames,
	})
	utils.DieIf(err)

	utils.Log("created", hostId)

	if optStatus != "" {
		err := newMackerel().UpdateHostStatus(hostId, optStatus)
		utils.DieIf(err)
	}
}

func doUpdate(c *cli.Context) {
	argHostId := c.Args().Get(0)
	optName := c.String("name")
	optStatus := c.String("status")
	optRoleFullnames := c.StringSlice("roleFullname")

	if argHostId = LoadHostIdFromConfig(); argHostId == "" {
		cli.ShowCommandHelp(c, "update")
		os.Exit(1)
	}

	isUpdated := false

	if optStatus != "" {
		err := newMackerel().UpdateHostStatus(argHostId, optStatus)
		utils.DieIf(err)

		isUpdated = true
	}
	if optName != "" || len(optRoleFullnames) > 0 {
		_, err := newMackerel().UpdateHost(argHostId, &mkr.UpdateHostParam{
			Name:          optName,
			RoleFullnames: optRoleFullnames,
		})
		utils.DieIf(err)

		isUpdated = true
	}

	if !isUpdated {
		cli.ShowCommandHelp(c, "update")
		os.Exit(1)
	}

	utils.Log("updated", argHostId)
}

func doThrow(c *cli.Context) {
	optHostId := c.String("host")
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
			utils.Log("warning", fmt.Sprintf("Failed to parse values: %s", err))
			continue
		}
		time, err := strconv.ParseInt(items[2], 10, 64)
		if err != nil {
			utils.Log("warning", fmt.Sprintf("Failed to parse values: %s", err))
			continue
		}

		metricValue := &mkr.MetricValue{
			Name:  items[0],
			Value: value,
			Time:  time,
		}

		metricValues = append(metricValues, metricValue)
	}
	utils.ErrorIf(scanner.Err())

	if optHostId != "" {
		err := newMackerel().PostHostMetricValuesByHostId(optHostId, metricValues)
		utils.DieIf(err)

		for _, metric := range metricValues {
			utils.Log("thrown", fmt.Sprintf("%s '%s\t%f\t%d'", optHostId, metric.Name, metric.Value, metric.Time))
		}
	} else if optService != "" {
		err := newMackerel().PostServiceMetricValues(optService, metricValues)
		utils.DieIf(err)

		for _, metric := range metricValues {
			utils.Log("thrown", fmt.Sprintf("%s '%s\t%f\t%d'", optService, metric.Name, metric.Value, metric.Time))
		}
	} else {
		cli.ShowCommandHelp(c, "throw")
		os.Exit(1)
	}
}

func doFetch(c *cli.Context) {
	argHostIds := c.Args()
	optMetricNames := c.StringSlice("name")

	if len(argHostIds) < 1 || len(optMetricNames) < 1 {
		cli.ShowCommandHelp(c, "fetch")
		os.Exit(1)
	}

	metricValues, err := newMackerel().FetchLatestMetricValues(argHostIds, optMetricNames)
	utils.DieIf(err)

	PrettyPrintJson(metricValues)
}

func doRetire(c *cli.Context) {
	argHostId := c.Args().Get(0)

	if argHostId = LoadHostIdFromConfig(); argHostId == "" {
		cli.ShowCommandHelp(c, "retire")
		os.Exit(1)
	}

	err := newMackerel().RetireHost(argHostId)
	utils.DieIf(err)
	utils.Log("retired", argHostId)
}
