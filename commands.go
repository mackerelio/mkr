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
	Usage: "Show host status",
	Description: `
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
	Usage: "Update host information like hostname, status and role",
	Description: `
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
`,
	Action: doThrow,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host, H", Value: "", Usage: "Post host metric values to <hostId>."},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Post service metric values to <service>."},
	},
}

var commandFetch = cli.Command{
	Name:  "fetch",
	Usage: "Fetch metric values",
	Description: `
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

func doStatus(c *cli.Context) {
	argHostId := c.Args().Get(0)
	isVerbose := c.Bool("verbose")

	if argHostId == "" {
		cli.ShowCommandHelp(c, "status")
		os.Exit(1)
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
	argName := c.String("name")
	argService := c.String("service")
	argRoles := c.StringSlice("role")
	argStatuses := c.StringSlice("status")

	hosts, err := newMackerel().FindHosts(&mkr.FindHostsParam{
		Name:     argName,
		Service:  argService,
		Roles:    argRoles,
		Statuses: argStatuses,
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
	argRoleFullnames := c.StringSlice("roleFullname")
	argStatus := c.String("status")

	if argHostName == "" {
		cli.ShowCommandHelp(c, "create")
		os.Exit(1)
	}

	hostId, err := newMackerel().CreateHost(&mkr.CreateHostParam{
		Name:          argHostName,
		RoleFullnames: argRoleFullnames,
	})
	utils.DieIf(err)

	if argStatus != "" {
		err := newMackerel().UpdateHostStatus(hostId, argStatus)
		utils.DieIf(err)
	}
}

func doUpdate(c *cli.Context) {
	argHostId := c.Args().Get(0)
	name := c.String("name")
	status := c.String("status")
	RoleFullnames := c.StringSlice("roleFullname")

	if argHostId == "" {
		cli.ShowCommandHelp(c, "update")
		os.Exit(1)
	}

	isUpdated := false

	if status != "" {
		err := newMackerel().UpdateHostStatus(argHostId, status)
		utils.DieIf(err)

		isUpdated = true
	}
	if name != "" || len(RoleFullnames) > 0 {
		_, err := newMackerel().UpdateHost(argHostId, &mkr.UpdateHostParam{
			Name:          name,
			RoleFullnames: RoleFullnames,
		})
		utils.DieIf(err)

		isUpdated = true
	}

	if !isUpdated {
		cli.ShowCommandHelp(c, "update")
		os.Exit(1)
	}
}

func doThrow(c *cli.Context) {
	argHostId := c.String("host")
	argService := c.String("service")

	var metricValues []*(mkr.MetricValue)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// name, value, timestamp
		// ex.) tcp.CLOSING 0 1397031808
		items := strings.Fields(line)
		fmt.Printf("%v+", items)
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

	if argHostId != "" {
		err := newMackerel().PostHostMetricValuesByHostId(argHostId, metricValues)
		utils.DieIf(err)
	} else if argService != "" {
		err := newMackerel().PostServiceMetricValues(argService, metricValues)
		utils.DieIf(err)
	} else {
		cli.ShowCommandHelp(c, "throw")
		os.Exit(1)
	}
}

func doFetch(c *cli.Context) {
	argHostIds := c.Args()
	argMetricNames := c.StringSlice("name")

	if len(argHostIds) < 1 || len(argMetricNames) < 1 {
		cli.ShowCommandHelp(c, "fetch")
		os.Exit(1)
	}

	metricValues, err := newMackerel().FetchLatestMetricValues(argHostIds, argMetricNames)
	utils.DieIf(err)

	PrettyPrintJson(metricValues)
}

func doRetire(c *cli.Context) {
	argHostId := c.Args().Get(0)

	if argHostId == "" {
		cli.ShowCommandHelp(c, "retire")
		os.Exit(1)
	}

	err := newMackerel().RetireHost(argHostId)
	utils.DieIf(err)
}
