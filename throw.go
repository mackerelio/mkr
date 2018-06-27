package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jpillora/backoff"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"gopkg.in/urfave/cli.v1"
)

var commandThrow = cli.Command{
	Name:      "throw",
	Usage:     "Post metric values",
	ArgsUsage: "[--host | -H <hostId>] [--service | -s <service>] stdin",
	Description: `
    Post metric values to 'host metric' or 'service metric'.
    Output format of metric values are compatible with that of a Sensu plugin.
    Requests "POST /api/v0/tsdb". See https://mackerel.io/api-docs/entry/host-metrics#post .
`,
	Action: doThrow,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host, H", Value: "", Usage: "Post host metric values to <hostID>."},
		cli.StringFlag{Name: "service, s", Value: "", Usage: "Post service metric values to <service>."},
	},
}

func doThrow(c *cli.Context) error {
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

		name := items[0]
		if optHostID != "" && !strings.HasPrefix(name, "custom.") {
			name = "custom." + name
		}

		metricValue := &mkr.MetricValue{
			Name:  name,
			Value: value,
			Time:  time,
		}

		metricValues = append(metricValues, metricValue)
	}
	logger.ErrorIf(scanner.Err())

	client := newMackerelFromContext(c)

	if optHostID != "" {
		logger.DieIf(requestWithRetry(func() error {
			return client.PostHostMetricValuesByHostID(optHostID, metricValues)
		}))

		for _, metric := range metricValues {
			logger.Log("thrown", fmt.Sprintf("%s '%s\t%f\t%d'", optHostID, metric.Name, metric.Value, metric.Time))
		}
	} else if optService != "" {
		logger.DieIf(requestWithRetry(func() error {
			return client.PostServiceMetricValues(optService, metricValues)
		}))

		for _, metric := range metricValues {
			logger.Log("thrown", fmt.Sprintf("%s '%s\t%f\t%d'", optService, metric.Name, metric.Value, metric.Time))
		}
	} else {
		cli.ShowCommandHelp(c, "throw")
		os.Exit(1)
	}
	return nil
}

const maxRetry = 10

func requestWithRetry(f func() error) error {
	b := &backoff.Backoff{
		Min:    1 * time.Second,
		Max:    1 * time.Minute,
		Factor: 2,
		Jitter: false,
	}
	var err error
	var delay time.Duration
	for int(b.Attempt()) < maxRetry {
		if b.Attempt() > 0 {
			logger.Log("warning", fmt.Sprintf("Failed to request. will retry after %.0f seconds. Error: %s", delay.Seconds(), err))
			time.Sleep(delay)
		}
		if err = f(); err == nil {
			// SUCCESS!!
			break
		} else if e, isAPIError := err.(*mkr.APIError); isAPIError {
			// Do not retry when status is 4XX
			if e.StatusCode >= 400 && e.StatusCode < 500 {
				break
			}
		}

		delay = b.Duration()
	}
	return err
}
