package main

import (
	"os"

	"github.com/codegangsta/cli"
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

	alerts, err := newMackerel(conffile).FindAlerts()
	logger.DieIf(err)

	PrettyPrintJSON(alerts)
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
