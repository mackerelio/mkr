package main

import (
	"os"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var commandChannels = cli.Command{
	Name:  "channels",
	Usage: "List notification channels",
	Description: `
	Lists notification channels.
	Requests APIs under "/api/v0/channels". See https://mackerel.io/api-docs/entry/channels .
	`,
	Action: doChannelsList,
}

func doChannelsList(c *cli.Context) error {
	// Waiting for mackerel-client-go to be bumped to version supporting FindChannels.
	client := mackerelclient.NewFromContext(c)
	channels, err := client.FindChannels()
	logger.DieIf(err)

	format.PrettyPrintJSON(os.Stdout, channels)
	return nil
}
