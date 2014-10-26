package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandHost,
	commandMetric,
}

var commandHost = cli.Command{
	Name:  "host",
	Usage: "",
	Description: `
`,
	Action: doHost,
}

var commandMetric = cli.Command{
	Name:  "metric",
	Usage: "",
	Description: `
`,
	Action: doMetric,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doHost(c *cli.Context) {
}

func doMetric(c *cli.Context) {
}
