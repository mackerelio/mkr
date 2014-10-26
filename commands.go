package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandStatus,
	commandHosts,
	commandCreate,
	commandUpdate,
	commandThrow,
	commandFetch,
	commandRetire,
}

var commandStatus = cli.Command{
	Name:  "status",
	Usage: "",
	Description: `
`,
	Action: doStatus,
}

var commandHosts = cli.Command{
	Name:  "hosts",
	Usage: "",
	Description: `
`,
	Action: doHosts,
}

var commandCreate = cli.Command{
	Name:  "create",
	Usage: "",
	Description: `
`,
	Action: doCreate,
}

var commandUpdate = cli.Command{
	Name:  "update",
	Usage: "",
	Description: `
`,
	Action: doUpdate,
}

var commandThrow = cli.Command{
	Name:  "throw",
	Usage: "",
	Description: `
`,
	Action: doThrow,
}

var commandFetch = cli.Command{
	Name:  "fetch",
	Usage: "",
	Description: `
`,
	Action: doFetch,
}

var commandRetire = cli.Command{
	Name:  "retire",
	Usage: "",
	Description: `
`,
	Action: doRetire,
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

func doStatus(c *cli.Context) {
}

func doHosts(c *cli.Context) {
}

func doCreate(c *cli.Context) {
}

func doUpdate(c *cli.Context) {
}

func doThrow(c *cli.Context) {
}

func doFetch(c *cli.Context) {
}

func doRetire(c *cli.Context) {
}
