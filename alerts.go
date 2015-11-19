package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codegangsta/cli"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
)

var commandAlerts = cli.Command{
	Name:  "alerts",
	Usage: "Retrieve/Close alerts",
	Description: `
    Retrieve/Close alerts. Without subcommand, show all alerts.
    Request APIs under "/api/v0/alerts". See http://help-ja.mackerel.io/entry/spec/api/v0 .
`,
	Action: doAlertsList,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "format, f", Value: "", Usage: "Output format. (human/json)"},
	},
	Subcommands: []cli.Command{
		{
			Name:        "close",
			Usage:       "close alerts",
			Description: "Pull monitor rules from Mackerel server and save them to a file. The file can be specified by filepath argument <file>. The default is 'monitors.json'.",
			Action:      doAlertsClose,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "reason, r", Value: "", Usage: "Reason of closing alert."},
				cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
			},
		},
	},
}

func doAlertsList(c *cli.Context) {
	conffile := c.GlobalString("conf")
	client := newMackerel(conffile)

	alerts, err := client.FindAlerts()
	logger.DieIf(err)
	if c.String("format") == "json" {
		PrettyPrintJSON(alerts)
		return
	}

	hostsJSON, err := client.FindHosts(&mkr.FindHostsParam{
		Statuses: []string{"working", "standby", "poweroff", "maintenance"},
	})
	hosts := map[string]*HostFormat{}
	for _, host := range hostsJSON {
		format := &HostFormat{
			ID:            host.ID,
			Name:          host.Name,
			Status:        host.Status,
			RoleFullnames: host.GetRoleFullnames(),
			IsRetired:     host.IsRetired,
			CreatedAt:     host.DateStringFromCreatedAt(),
			IPAddresses:   host.IPAddresses(),
		}
		hosts[host.ID] = format
	}

	monitorsJSON, err := client.FindMonitors()
	monitors := map[string]*mkr.Monitor{}
	for _, monitor := range monitorsJSON {
		monitors[monitor.ID] = monitor
	}

	const layout = "2006-01-02 15:04:05"
	for _, alert := range alerts {
		host := hosts[alert.HostID]
		hostMsg := ""
		if host != nil {
			hostMsg = fmt.Sprintf("%s %s %s", host.Name, host.Status, host.RoleFullnames)
		}
		monitor := monitors[alert.MonitorID]
		monitorMsg := ""
		if monitor != nil {
			switch monitor.Type {
			case "connectivity":
				monitorMsg = fmt.Sprintf("%s", monitor.Type)
			case "host":
				switch alert.Status {
				case "CRITICAL":
					monitorMsg = fmt.Sprintf("%s %.2f %s %.2f", monitor.Metric, alert.Value, monitor.Operator, monitor.Critical)
				case "WARNING":
					monitorMsg = fmt.Sprintf("%s %.2f %s %.2f", monitor.Metric, alert.Value, monitor.Operator, monitor.Warning)
				default:
					monitorMsg = fmt.Sprintf("%s %.2f %s %.2f", monitor.Metric, alert.Value, monitor.Operator, monitor.Critical)
				}
			default:
				monitorMsg = fmt.Sprintf("%s", monitor.Type)
			}
		}
		fmt.Printf("%s %s %8s %s %s\n", alert.ID, time.Unix(alert.OpenedAt, 0).Format(layout), alert.Status, monitorMsg, hostMsg)
	}
}

func doAlertsClose(c *cli.Context) {
	conffile := c.GlobalString("conf")
	isVerbose := c.Bool("verbose")
	argAlertIDs := c.Args()
	reason := c.String("reason")

	if len(argAlertIDs) < 1 {
		cli.ShowCommandHelp(c, "alerts")
		os.Exit(1)
	}

	client := newMackerel(conffile)
	for _, alertID := range argAlertIDs {
		alert, err := client.CloseAlert(alertID, reason)
		logger.DieIf(err)

		logger.Log("Alert closed", alertID)
		if isVerbose == true {
			PrettyPrintJSON(alert)
		}
	}
}
