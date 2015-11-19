package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
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
		cli.BoolFlag{Name: "color, c", Usage: "Colorized output. default: true"},
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

type alertSet struct {
	Alert   *mkr.Alert
	Host    *mkr.Host
	Monitor *mkr.Monitor
}

func joinMonitorsAndHosts(client *mkr.Client, alerts []*mkr.Alert) []*alertSet {
	hostsJSON, err := client.FindHosts(&mkr.FindHostsParam{
		Statuses: []string{"working", "standby", "poweroff", "maintenance"},
	})
	logger.DieIf(err)

	hosts := map[string]*mkr.Host{}
	for _, host := range hostsJSON {
		hosts[host.ID] = host
	}

	monitorsJSON, err := client.FindMonitors()
	logger.DieIf(err)

	monitors := map[string]*mkr.Monitor{}
	for _, monitor := range monitorsJSON {
		monitors[monitor.ID] = monitor
	}

	alertSets := []*alertSet{}
	for _, alert := range alerts {
		alertSets = append(
			alertSets,
			&alertSet{Alert: alert, Host: hosts[alert.HostID], Monitor: monitors[alert.MonitorID]},
		)
	}
	return alertSets
}

func formatJoinedAlert(alertSet *alertSet, colorize bool) string {
	const layout = "2006-01-02 15:04:05"

	host := alertSet.Host
	monitor := alertSet.Monitor
	alert := alertSet.Alert

	hostMsg := ""
	if host != nil {
		statusMsg := host.Status
		if host.IsRetired == true {
			statusMsg = "retired"
		}
		if colorize {
			switch statusMsg {
			case "working":
				statusMsg = color.BlueString("working")
			case "standby":
				statusMsg = color.GreenString("standby")
			case "poweroff":
				statusMsg = "poweroff"
			case "maintenance":
				statusMsg = color.YellowString("maintenance")
			}
		}
		hostMsg = fmt.Sprintf("%s %s", host.Name, statusMsg)
		roleMsgs := []string{}
		for service, roles := range host.Roles {
			roleMsgs = append(roleMsgs, fmt.Sprintf("%s:%s", service, strings.Join(roles, ",")))
		}
		hostMsg += " [" + strings.Join(roleMsgs, ", ") + "]"
	}

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
	statusMsg := alert.Status
	if colorize {
		switch alert.Status {
		case "CRITICAL":
			statusMsg = color.RedString("CRITICAL")
		case "WARNING":
			statusMsg = color.YellowString("WARNING ")
		}
	}
	return fmt.Sprintf("%s %s %s %s %s", alert.ID, time.Unix(alert.OpenedAt, 0).Format(layout), statusMsg, monitorMsg, hostMsg)
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

	joinedAlerts := joinMonitorsAndHosts(client, alerts)

	for _, joinAlert := range joinedAlerts {
		fmt.Println(formatJoinedAlert(joinAlert, c.Bool("color")))
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
