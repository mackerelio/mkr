package alerts

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:      "alerts",
	Usage:     "Retrieve/Close alerts",
	ArgsUsage: "[--with-closed | -w] [--limit | -l] [--jq <formula>]",
	Description: `
    Retrieve/Close alerts. With no subcommand specified, this will show all alerts.
    Requests APIs under "/api/v0/alerts". See https://mackerel.io/api-docs/entry/alerts .
`,
	Action: doAlertsRetrieve,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "with-closed, w", Usage: "Display open alert including close alert. default: false"},
		cli.IntFlag{Name: "limit, l", Value: defaultAlertsLimit, Usage: fmt.Sprintf("Set the number of alerts to display. Default is set to %d when -with-closed is set, otherwise all the open alerts are displayed.", defaultAlertsLimit)},
		cli.StringFlag{Name: "jq", Usage: "Query to select values from the response using jq syntax"},
	},
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "list alerts",
			ArgsUsage: "[--service | -s <service>] [--host-status | -S <file>] [--color | -c] [--with-closed | -w] [--limit | -l]",
			Description: `
    Shows alerts in human-readable format.
`,
			Action: doAlertsList,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "service, s",
					Value: &cli.StringSlice{},
					Usage: "Filters alerts by service. Multiple choices are allowed.",
				},
				cli.StringSliceFlag{
					Name:  "host-status, S",
					Value: &cli.StringSlice{},
					Usage: "Filters alerts by status of each host. Multiple choices are allowed.",
				},
				cli.BoolTFlag{Name: "color, c", Usage: "Colorize output. default: true"},
				cli.BoolFlag{Name: "with-closed, w", Usage: "Display open alert including close alert. default: false"},
				cli.IntFlag{Name: "limit, l", Value: defaultAlertsLimit, Usage: fmt.Sprintf("Set the number of alerts to display. Default is set to %d when -with-closed is set, otherwise all the open alerts are displayed.", defaultAlertsLimit)},
			},
		},
		{
			Name:      "close",
			Usage:     "close alerts",
			ArgsUsage: "<alertIds....>",
			Description: `
    Closes alerts. Multiple alert IDs can be specified.
`,
			Action: doAlertsClose,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "reason, r", Value: "", Usage: "Reason of closing alert."},
				cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
			},
		},
	},
}

const defaultAlertsLimit int = 100

type alertSet struct {
	Alert   *mackerel.Alert
	Host    *mackerel.Host
	Monitor mackerel.Monitor
}

func joinMonitorsAndHosts(client *mackerel.Client, alerts []*mackerel.Alert) []*alertSet {
	hostsJSON, err := client.FindHosts(&mackerel.FindHostsParam{
		Statuses: []string{"working", "standby", "poweroff", "maintenance"},
	})
	logger.DieIf(err)

	hosts := map[string]*mackerel.Host{}
	for _, host := range hostsJSON {
		hosts[host.ID] = host
	}

	monitorsJSON, err := client.FindMonitors()
	logger.DieIf(err)

	monitors := map[string]mackerel.Monitor{}
	for _, monitor := range monitorsJSON {
		monitors[monitor.MonitorID()] = monitor
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
		if host.IsRetired {
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
		hostMsg = fmt.Sprintf(" %s %s", host.Name, statusMsg)
		roleMsgs := []string{}
		for service, roles := range host.Roles {
			roleMsgs = append(roleMsgs, fmt.Sprintf("%s:%s", service, strings.Join(roles, ",")))
		}
		hostMsg += " [" + strings.Join(roleMsgs, ", ") + "]"
	}

	monitorMsg := ""
	if monitor != nil {
		switch m := monitor.(type) {
		case *mackerel.MonitorConnectivity:
			monitorMsg = ""
		case *mackerel.MonitorHostMetric:
			if alert.Status == "CRITICAL" && m.Critical != nil {
				monitorMsg = fmt.Sprintf("%s %.2f %s %.2f", m.Metric, alert.Value, m.Operator, *m.Critical)
			} else if alert.Status == "WARNING" && m.Warning != nil {
				monitorMsg = fmt.Sprintf("%s %.2f %s %.2f", m.Metric, alert.Value, m.Operator, *m.Warning)
			} else {
				monitorMsg = fmt.Sprintf("%s %.2f", m.Metric, alert.Value)
			}
		case *mackerel.MonitorServiceMetric:
			if alert.Status == "CRITICAL" && m.Critical != nil {
				monitorMsg = fmt.Sprintf("%s %s %.2f %s %.2f", m.Service, m.Metric, alert.Value, m.Operator, *m.Critical)
			} else if alert.Status == "WARNING" && m.Warning != nil {
				monitorMsg = fmt.Sprintf("%s %s %.2f %s %.2f", m.Service, m.Metric, alert.Value, m.Operator, *m.Warning)
			} else {
				monitorMsg = fmt.Sprintf("%s %s %.2f", m.Service, m.Metric, alert.Value)
			}
		case *mackerel.MonitorExternalHTTP:
			statusRegexp, _ := regexp.Compile("^[2345][0-9][0-9]$")
			switch alert.Status {
			case "CRITICAL":
				if statusRegexp.MatchString(alert.Message) && m.ResponseTimeCritical != nil {
					monitorMsg = fmt.Sprintf("%s %.2f > %.2f msec, status:%s", m.URL, alert.Value, *m.ResponseTimeCritical, alert.Message)
				} else {
					monitorMsg = fmt.Sprintf("%s %.2f msec, %s", m.URL, alert.Value, alert.Message)
				}
			case "WARNING":
				if statusRegexp.MatchString(alert.Message) && m.ResponseTimeWarning != nil {
					monitorMsg = fmt.Sprintf("%.2f > %.2f msec, status:%s", alert.Value, *m.ResponseTimeWarning, alert.Message)
				} else {
					monitorMsg = fmt.Sprintf("%.2f msec, %s", alert.Value, alert.Message)
				}
			default:
				monitorMsg = fmt.Sprintf("%.2f msec, status:%s", alert.Value, alert.Message)
			}
		case *mackerel.MonitorExpression:
			expression := formatExpressionOneline(m.Expression)
			if alert.Status == "CRITICAL" && m.Critical != nil {
				monitorMsg = fmt.Sprintf("%s %.2f %s %.2f", expression, alert.Value, m.Operator, *m.Critical)
			} else if alert.Status == "WARNING" && m.Warning != nil {
				monitorMsg = fmt.Sprintf("%s %.2f %s %.2f", expression, alert.Value, m.Operator, *m.Warning)
			} else if alert.Status == "UNKNOWN" {
				monitorMsg = expression
			} else {
				monitorMsg = fmt.Sprintf("%s %.2f", expression, alert.Value)
			}
		case *mackerel.MonitorAnomalyDetection:
			monitorMsg = ""
		default:
			monitorMsg = monitor.MonitorType()
		}
		if monitorMsg == "" {
			monitorMsg = monitor.MonitorName()
		} else {
			monitorMsg = monitor.MonitorName() + " " + monitorMsg
		}
	}
	// If alert is caused by check monitoring, take monitorMsg from alert.message
	if alert.Type == "check" {
		monitorMsg = formatCheckMessage(alert.Message)
	}

	statusMsg := alert.Status
	if colorize {
		switch alert.Status {
		case "CRITICAL":
			statusMsg = color.RedString("CRITICAL ")
		case "WARNING":
			statusMsg = color.YellowString("WARNING ")
		case "OK":
			statusMsg = color.GreenString("OK ")
		case "UNKNOWN":
			statusMsg = "UNKNOWN "
		}
	}
	return fmt.Sprintf("%s %s %s %s%s", alert.ID, time.Unix(alert.OpenedAt, 0).Format(layout), statusMsg, monitorMsg, hostMsg)
}

var expressionNewlinePattern = regexp.MustCompile(`\s*[\r\n]+\s*`)

func formatExpressionOneline(expr string) string {
	expr = strings.Trim(expressionNewlinePattern.ReplaceAllString(expr, " "), " ")
	return strings.Replace(strings.Replace(expr, "( ", "(", -1), " )", ")", -1)
}

func formatCheckMessage(msg string) string {
	truncated := false
	if index := strings.IndexAny(msg, "\n\r"); index != -1 {
		msg = msg[0:index]
		truncated = true
	}
	if runes := []rune(msg); len(runes) > 100 {
		msg = string(runes[0:100])
		truncated = true
	}
	if truncated {
		msg = msg + "..."
	}
	return msg
}

func doAlertsRetrieve(c *cli.Context) error {
	client := mackerelclient.NewFromContext(c)
	withClosed := c.Bool("with-closed")
	alerts, err := fetchAlerts(client, withClosed, getAlertsLimit(c, withClosed))
	logger.DieIf(err)
	err = format.PrettyPrintJSON(os.Stdout, alerts, c.String("jq"))
	logger.DieIf(err)
	return nil
}

func doAlertsList(c *cli.Context) error {
	filterServices := c.StringSlice("service")
	filterStatuses := c.StringSlice("host-status")
	client := mackerelclient.NewFromContext(c)
	withClosed := c.Bool("with-closed")
	alerts, err := fetchAlerts(client, withClosed, getAlertsLimit(c, withClosed))
	logger.DieIf(err)

	joinedAlerts := joinMonitorsAndHosts(client, alerts)
	for _, joinAlert := range joinedAlerts {
		if len(filterServices) > 0 {
			found := false
			for _, filterService := range filterServices {
				if joinAlert.Host != nil {
					if _, ok := joinAlert.Host.Roles[filterService]; ok {
						found = true
					}
				} else {
					var service string
					if m, ok := joinAlert.Monitor.(*mackerel.MonitorServiceMetric); ok {
						service = m.Service
					} else if m, ok := joinAlert.Monitor.(*mackerel.MonitorExternalHTTP); ok {
						service = m.Service
					}
					found = service == filterService
				}
			}
			if !found {
				continue
			}
		}
		if len(filterStatuses) > 0 {
			found := false
			for _, filterStatus := range filterStatuses {
				if joinAlert.Host != nil && joinAlert.Host.Status == filterStatus {
					found = true
				}
			}
			if !found {
				continue
			}
		}
		fmt.Fprintln(color.Output, formatJoinedAlert(joinAlert, c.BoolT("color")))
	}
	return nil
}

func getAlertsLimit(c *cli.Context, withClosed bool) int {
	if c.IsSet("limit") {
		return c.Int("limit")
	}
	if withClosed {
		return defaultAlertsLimit
	}
	// When -limit is not set, mkr alerts should print all the open alerts.
	return math.MaxInt32
}

func fetchAlerts(client *mackerel.Client, withClosed bool, limit int) ([]*mackerel.Alert, error) {
	if limit < 0 {
		return nil, errors.New("limit should not be negative")
	}
	var resp *mackerel.AlertsResp
	var err error
	if withClosed {
		if resp, err = client.FindWithClosedAlerts(); err != nil {
			return nil, err
		}
		if resp.NextID != "" {
			for {
				if limit <= len(resp.Alerts) {
					break
				}
				nextResp, err := client.FindWithClosedAlertsByNextID(resp.NextID)
				if err != nil {
					return nil, err
				}
				resp.Alerts = append(resp.Alerts, nextResp.Alerts...)
				resp.NextID = nextResp.NextID
				if resp.NextID == "" {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
	} else {
		if resp, err = client.FindAlerts(); err != nil {
			return nil, err
		}
		if resp.NextID != "" {
			for {
				if limit <= len(resp.Alerts) {
					break
				}
				nextResp, err := client.FindAlertsByNextID(resp.NextID)
				if err != nil {
					return nil, err
				}
				resp.Alerts = append(resp.Alerts, nextResp.Alerts...)
				resp.NextID = nextResp.NextID
				if resp.NextID == "" {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
	if len(resp.Alerts) > limit {
		resp.Alerts = resp.Alerts[:limit]
	}
	return resp.Alerts, nil
}

func doAlertsClose(c *cli.Context) error {
	isVerbose := c.Bool("verbose")
	argAlertIDs := c.Args()
	reason := c.String("reason")

	if len(argAlertIDs) < 1 {
		cli.ShowCommandHelpAndExit(c, "alerts", 1)
	}

	client := mackerelclient.NewFromContext(c)
	for _, alertID := range argAlertIDs {
		alert, err := client.CloseAlert(alertID, reason)
		logger.DieIf(err)

		logger.Log("Alert closed", alertID)
		if isVerbose {
			err := format.PrettyPrintJSON(os.Stdout, alert, "")
			logger.DieIf(err)
		}
	}
	return nil
}
