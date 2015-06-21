package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/mackerelio/mackerel-agent/config"
)

func main() {
	app := cli.NewApp()
	app.Name = "mkr"
	app.Version = Version
	app.Usage = "A CLI tool for mackerel.io"
	app.Author = "Hatena Co., Ltd."
	app.Commands = Commands
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Value: config.DefaultConfig.Conffile,
			Usage: "Config file path",
		},
	}

	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	app.Run(os.Args)
}
