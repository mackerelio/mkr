package channels

import (
	"os"

	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v2"
)

// Command is the definition of channels subcommand
var Command = cli.Command{
	Name:  "channels",
	Usage: "List notification channels",
	Description: `
	Lists notification channels. With no subcommand specified, this will show all channels.
	Requests APIs under "/api/v0/channels". See https://mackerel.io/api-docs/entry/channels .
	`,
	Action: doChannels,
	Flags: []cli.Flag{
		jq.CommandLineFlag,
	},
	Subcommands: []cli.Command{
		{
			Name:      "pull",
			Usage:     "pull channel settings",
			ArgsUsage: "[--file-path | -F <file>] [--verbose | -v]",
			Description: `
    Pull channels settings from Mackerel server and save them to a file. The file can be specified by filepath argument <file>. The default is 'channels.json'.
`,
			Action: doChannelsPull,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file-path, F", Value: "", Usage: "Filename to store channel settings. default: channels.json"},
				cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
			},
		},
	},
}

func doChannels(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&channelsApp{
		client:    client,
		outStream: os.Stdout,
		jqFilter:  c.String("jq"),
	}).run()
}

func doChannelsPull(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	isVerbose := c.Bool("verbose")
	filePath := c.String("file-path")

	return (&channelsApp{
		client:    client,
		outStream: os.Stdout,
	}).pullChannels(isVerbose, filePath)
}
