package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mkr/logger"
	"github.com/urfave/cli/v3"
)

func main() {
	version, gitcommit := fromVCS()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app := &cli.Command{
		Name:    "mkr",
		Version: fmt.Sprintf("%s (rev:%s)", version, gitcommit),
		Usage:   "A CLI tool for mackerel.io",
		Authors: []any{
			string("Hatena Co., Ltd."),
		},
		Commands: Commands,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "conf",
				Value: config.DefaultConfig.Conffile,
				Usage: "Config `file` path",
			},
			&cli.StringFlag{
				Name: "apibase",
				// this default value is set in config.LoadApibaseFromConfigWithFallback
				Usage: fmt.Sprintf("API Base `URL` (default: \"%s\")", config.DefaultConfig.Apibase),
			},
		},
	}

	err := app.Run(ctx, os.Args)
	if err != nil {
		exitCode := 1
		if excoder, ok := err.(cli.ExitCoder); ok {
			exitCode = excoder.ExitCode()
		}
		logger.Log("error", err.Error())
		os.Exit(exitCode)
	}
}

func fromVCS() (version, rev string) {
	version = "unknown"
	rev = "unknown"
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	version = info.Main.Version
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			rev = s.Value
			return
		}
	}
	return
}
