package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mackerelio/mackerel-agent/config"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "mkr"
	app.Version = fmt.Sprintf("%s (rev:%s)", version, gitcommit)
	app.Usage = "A CLI tool for mackerel.io"
	app.Author = "Hatena Co., Ltd."
	app.Commands = Commands
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Value: config.DefaultConfig.Conffile,
			Usage: "Config file path",
		},
		cli.StringFlag{
			Name:  "apibase",
			Value: config.DefaultConfig.Apibase,
			Usage: "API Base",
		},
	}

	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	app.Run(os.Args)
}
