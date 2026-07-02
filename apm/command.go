package apm

import (
	"context"
	"os"
	"time"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v3"
)

// Command is the definition of apm subcommand
var Command = &cli.Command{
	Name:  "apm",
	Usage: "APM related commands",
	Commands: []*cli.Command{
		httpServerStatsCommand,
	},
}

// httpServerStatsCommand is the definition of http-server-stats subcommand
var httpServerStatsCommand = &cli.Command{
	Name:      "http-server-stats",
	Usage:     "List HTTP server stats of a service",
	ArgsUsage: "[--service <serviceName>]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      "service",
			Usage:     "Service name",
			Aliases:   []string{"s"},
			Required:  true,
			TakesFile: false,
		},
		&cli.StringFlag{
			Name:      "service-namespace",
			Usage:     "Service namespace",
			Aliases:   []string{},
			Required:  false,
			TakesFile: false,
		},
		&cli.StringFlag{
			Name:      "environment",
			Usage:     "Environment name",
			Aliases:   []string{},
			Required:  false,
			TakesFile: false,
		},
		&cli.StringFlag{
			Name:      "route",
			Usage:     "Filter by route",
			Aliases:   []string{},
			Required:  false,
			TakesFile: false,
		},
		&cli.TimestampFlag{
			Name:      "from",
			Usage:     "Start timestamp",
			Aliases:   []string{},
			Required:  false,
			TakesFile: false,
			Config: cli.TimestampConfig{
				Timezone: time.Local,
				Layouts:  []string{time.RFC3339, time.DateTime},
			},
			Value: time.Now().Truncate(time.Minute).Add(-30 * time.Minute),
		},
		&cli.TimestampFlag{
			Name:      "to",
			Usage:     "End timestamp",
			Aliases:   []string{},
			Required:  false,
			TakesFile: false,
			Config: cli.TimestampConfig{
				Timezone: time.Local,
				Layouts:  []string{time.RFC3339, time.DateTime},
			},
			Value: time.Now().Truncate(time.Minute),
		},
		&cli.IntFlag{
			Name:      "page",
			Usage:     "Page number",
			Aliases:   []string{},
			Required:  false,
			TakesFile: false,
			Value:     1,
		},
	},
	Action: doHTTPServerStats,
}

func doHTTPServerStats(ctx context.Context, c *cli.Command) error {
	client := mackerelclient.NewFromCliCommand(c)
	app := &httpServerStatsApp{
		client:    client,
		logger:    logger.New(),
		outStream: os.Stdout,
	}
	param := &mackerel.ListHTTPServerStatsParam{
		ServiceName: c.String("service"),
		From:        c.Timestamp("from").Unix(),
		To:          c.Timestamp("to").Unix(),
	}
	if ns := c.String("service-namespace"); ns != "" {
		param.ServiceNamespace = &ns
	}
	if env := c.String("environment"); env != "" {
		param.Environment = &env
	}
	if route := c.String("route"); route != "" {
		param.Route = &route
	}
	if page := c.Int("page"); page != 0 {
		param.Page = &page
	}
	return app.listHTTPServerStats(ctx, param)
}
