package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mackerelio/gomkr/utils"
	mkr "github.com/mackerelio/mackerel-client-go"
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
	Usage: "Show host status",
	Description: `
`,
	Action: doStatus,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "verbose, V", Usage: "Output verbose info"},
	},
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

func newMackerel() *mkr.Client {
	apiKey := os.Getenv("MACKEREL_APIKEY")
	if apiKey == "" {
		utils.Log("error", `
Not set MACKEREL_APIKEY environment variable. (Try "export MACKEREL_APIKEY='<Your apikey>'")
`)
		os.Exit(1)
	}
	mackerel := mkr.NewClient(apiKey)
	return mackerel
}

func doStatus(c *cli.Context) {
	argHostId := c.Args().Get(0)

	if argHostId == "" {
		cli.ShowCommandHelp(c, "status")
		os.Exit(1)
	}

	mackerel := newMackerel()
	host, err := mackerel.FindHost(argHostId)
	utils.DieIf(err)

	data, err := json.MarshalIndent(host, "", "    ")
	utils.DieIf(err)

	fmt.Fprintln(os.Stdout, string(data))
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
