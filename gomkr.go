package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gomkr"
	app.Version = Version
	app.Usage = "mackerel.io CLI tool"
	app.Author = "Hatena"
	app.Email = "y_uuki@hatena.ne.jp"
	app.Commands = Commands

	app.Run(os.Args)
}
