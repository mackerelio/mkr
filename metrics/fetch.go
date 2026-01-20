package metrics

import (
	"os"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v2"
)

var CommandFetch = cli.Command{
	Name:      "fetch",
	Usage:     "Fetch latest metric values",
	ArgsUsage: "[--name | -n <metricName>] [--jq <formula>] hostIds...",
	Description: `
    Fetch latest metric values about the hosts.
    Requests "GET /api/v0/tsdb/latest". See https://mackerel.io/api-docs/entry/host-metrics#get-latest .
`,
	Action: doFetch,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   &cli.StringSlice{},
			Usage:   "Fetch metric values identified with <name>. Required. Multiple choices are allowed. ",
		},
		jq.CommandLineFlag,
	},
}

func doFetch(c *cli.Context) error {
	argHostIDs := c.Args()
	optMetricNames := c.StringSlice("name")

	if len(argHostIDs) < 1 || len(optMetricNames) < 1 {
		cli.ShowCommandHelpAndExit(c, "fetch", 1)
	}

	allMetricValues := make(mackerel.LatestMetricValues)
	// Fetches 100 hosts per one request (to avoid URL maximum length).
	for _, hostIds := range split(argHostIDs, 100) {
		metricValues, err := mackerelclient.NewFromContext(c).FetchLatestMetricValues(hostIds, optMetricNames)
		logger.DieIf(err)
		for key := range metricValues {
			allMetricValues[key] = metricValues[key]
		}
	}

	err := format.PrettyPrintJSON(os.Stdout, allMetricValues, c.String("jq"))
	logger.DieIf(err)
	return nil
}

func split(ids []string, count int) [][]string {
	xs := make([][]string, 0, (len(ids)+count-1)/count)
	for i, name := range ids {
		if i%count == 0 {
			xs = append(xs, []string{})
		}
		xs[len(xs)-1] = append(xs[len(xs)-1], name)
	}
	return xs
}
