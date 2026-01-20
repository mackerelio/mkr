package wrap

import (
	"fmt"
	"os"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mkr/logger"
	"github.com/urfave/cli/v2"
)

// Command is definition of mkr wrap
var Command = &cli.Command{
	Name:      "wrap",
	Usage:     "Wrap and monitor batch jobs to run with cron etc",
	ArgsUsage: "[--name|-n <name>] [OPTIONS] -- /path/to/batch",
	Description: `
    Wrap a batch command with specifying it as arguments. If the command failed
    with non-zero exit code, it sends a report to Mackerel and raises an alert.
    It is useful for cron jobs etc.
`,
	Action: doWrap,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   "",
			Usage:   "The `check-name` which must be unique on a host. If it is empty it will be automatically derived.",
		},
		&cli.BoolFlag{
			Name:    "detail",
			Aliases: []string{"d"},
			Usage:   "send a detailed report contains command output",
		},
		&cli.StringFlag{
			Name:    "note",
			Aliases: []string{"N"},
			Value:   "",
			Usage:   "`note` of the job",
		},
		&cli.StringFlag{
			Name:    "host",
			Aliases: []string{"H"},
			Value:   "",
			Usage:   "`hostID`",
		},
		&cli.BoolFlag{
			Name:    "warning",
			Aliases: []string{"w"},
			Usage:   "alerts as warning",
		},
		&cli.BoolFlag{
			Name:    "auto-close",
			Aliases: []string{"a"},
			Usage:   "automatically close an existing alert when the command success",
		},
		&cli.DurationFlag{
			Name:    "notification-interval",
			Aliases: []string{"I"},
			Usage:   "The notification re-sending `interval`. If it is zero, never re-send. (minimum 10 minutes)",
		},
		// XXX Implementation of maxCheckAttempts is difficult because the
		// execution interval of cron or batches are not always one-minute.
		// This is due to the server-side logic of the Mackerel.
	},
}

func doWrap(c *cli.Context) error {
	confFile := c.String("conf")
	var conf *config.Config
	if _, err := os.Stat(confFile); err == nil {
		conf, err = config.LoadConfig(confFile)
		if err != nil {
			logger.Logf("error", "[mkr wrap] failed to load the config %q: %s", confFile, err)
		}
	} else {
		logger.Logf("info", "[mkr wrap] configuraion file %q not found", confFile)
	}
	if conf == nil {
		// fallback default config
		conf = config.DefaultConfig
	}

	apibase := c.String("apibase")
	if apibase == "" {
		apibase = conf.Apibase
	}

	apikey := os.Getenv("MACKEREL_APIKEY")
	if apikey == "" {
		apikey = conf.Apikey
	}
	if apikey == "" {
		logger.Log("error", "[mkr wrap] failed to detect Mackerel APIKey. Try to specify in mackerel-agent.conf or export MACKEREL_APIKEY='<Your apikey>'")
	}
	var hostID string
	if id := c.String("host"); id != "" {
		hostID = id
	} else {
		hostID, _ = conf.LoadHostID()
	}
	if hostID == "" {
		logger.Log("error", "[mkr wrap] failed to load hostID. Try to specify -host option explicitly")
	}
	// Since command execution has the highest priority, even when the config
	// loading is failed, or apikey or hostID is empty, we don't return errors
	// and only output the log here.

	cmd := c.Args().Slice()
	if len(cmd) > 0 && cmd[0] == "--" {
		cmd = cmd[1:]
	}
	if len(cmd) < 1 {
		return fmt.Errorf("no commands specified")
	}

	return (&wrap{
		apibase:              apibase,
		name:                 c.String("name"),
		detail:               c.Bool("detail"),
		note:                 c.String("note"),
		warning:              c.Bool("warning"),
		autoClose:            c.Bool("auto-close"),
		notificationInterval: c.Duration("notification-interval"),
		hostID:               hostID,
		apikey:               apikey,
		cmd:                  cmd,
		outStream:            os.Stdout,
		errStream:            os.Stderr,
	}).run()
}
