package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var commandChannels = cli.Command{
	Name:  "channels",
	Usage: "List notification channels",
	Description: `
	Lists notification channels. With no subcommand specified, this will show all channels.
	Requests APIs under "/api/v0/channels". See https://mackerel.io/api-docs/entry/channels .
	`,
	Action: doChannelsList,
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

func doChannelsList(c *cli.Context) error {
	client := mackerelclient.NewFromContext(c)
	channels, err := client.FindChannels()
	logger.DieIf(err)

	format.PrettyPrintJSON(os.Stdout, channels)
	return nil
}

func doChannelsPull(c *cli.Context) error {
	isVerbose := c.Bool("verbose")
	filePath := c.String("file-path")

	channels, err := mackerelclient.NewFromContext(c).FindChannels()
	logger.DieIf(err)

	channelSaveRules(channels, filePath)

	if isVerbose {
		format.PrettyPrintJSON(os.Stdout, channels)
	}

	if filePath == "" {
		filePath = "channels.json"
	}
	logger.Log("info", fmt.Sprintf("Channels are saved to '%s' (%d rules).", filePath, len(channels)))
	return nil
}

func channelSaveRules(rules []*mackerel.Channel, optFilePath string) error {
	filePath := "channels.json"
	if optFilePath != "" {
		filePath = optFilePath
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	channels := map[string]interface{}{"channels": rules}
	data := format.JSONMarshalIndent(channels, "", "    ") + "\n"

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}
