package channels

import (
	"os"

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
	Action: doChannels,
}

func doChannels(c *cli.Context) error {
	client, err := mackerelclient.New(c.GlobalString("conf"), c.GlobalString("apibase"))
	if err != nil {
		return err
	}

	return (&channelsApp{
		client:    client,
		outStream: os.Stdout,
	}).run()
}
