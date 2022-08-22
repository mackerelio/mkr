package metrics

import (
	"os"
	"time"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:      "metrics",
	Usage:     "Fetch metric values",
	ArgsUsage: "[--host | -H <hostId>] [--service | -s <service>] [--name | -n <metricName>] [--jq <formula>] --from int --to int",
	Description: `
    Fetch metric values of 'host metric' or 'service metric'.
    Requests "/api/v0/hosts/<hostId>/metrics" or "/api/v0/services/<serviceName>/tsdb".
    See https://mackerel.io/api-docs/entry/host-metrics#get, https://mackerel.io/api-docs/entry/service-metrics#get.
`,
	Action: doMetrics,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host, H", Value: "", Usage: "Fetch host metric values of <hostID>."},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Fetch service metric values of <service>."},
		cli.StringFlag{Name: "name, n", Value: "", Usage: "The name of the metric for which you want to obtain the metric."},
		cli.Int64Flag{Name: "from", Usage: "The first of the period for which you want to obtain the metric. (epoch seconds)"},
		cli.Int64Flag{Name: "to", Usage: "The end of the period for which you want to obtain the metric. (epoch seconds)"},
		jq.CommandLineFlag,
	},
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
	jq := c.String("jq")

	client := mackerelclient.NewFromContext(c)

	if optHostID != "" {
		metricValue, err := client.FetchHostMetricValues(optHostID, optMetricName, from, to)
		logger.DieIf(err)

		err = format.PrettyPrintJSON(os.Stdout, metricValue, jq)
		logger.DieIf(err)
	} else if optService != "" {
		metricValue, err := client.FetchServiceMetricValues(optService, optMetricName, from, to)
		logger.DieIf(err)

		err = format.PrettyPrintJSON(os.Stdout, metricValue, jq)
		logger.DieIf(err)
	} else {
		cli.ShowCommandHelpAndExit(c, "metrics", 1)
	}
	return nil
}
