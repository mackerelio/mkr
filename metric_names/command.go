package metric_names

import (
	"os"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:      "metric-names",
	Usage:     "Fetch metric names",
	ArgsUsage: "[--host | -H <hostId>] [--service | -s <service>] [--jq <formula>]",
	Description: `
    Fetch metric names of 'host metric' or 'service metric'.
    Requests "/api/v0/hosts/<hostId>/metric-names" or "/api/v0/services/<serviceName>/metric-names".
    See https://mackerel.io/ja/api-docs/entry/hosts#metric-names, https://mackerel.io/ja/api-docs/entry/services#metric-names
`,
	Action: doMetricNames,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host, H", Value: "", Usage: "Fetch host metric names of <hostID>."},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Fetch service metric names of <service>."},
		jq.CommandLineFlag,
	},
}

func doMetricNames(c *cli.Context) error {
	optHostID := c.String("host")
	optService := c.String("service")
	jq := c.String("jq")

	client := mackerelclient.NewFromContext(c)

	if optHostID != "" {
		metricNames, err := client.ListHostMetricNames(optHostID)
		logger.DieIf(err)

		err = format.PrettyPrintJSON(os.Stdout, metricNames, jq)
		logger.DieIf(err)
	} else if optService != "" {
		metricNames, err := client.ListServiceMetricNames(optService)
		logger.DieIf(err)

		err = format.PrettyPrintJSON(os.Stdout, metricNames, jq)
		logger.DieIf(err)
	} else {
		cli.ShowCommandHelpAndExit(c, "metric-names", 1)
	}
	return nil
}
