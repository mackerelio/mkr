package channels

import (
	"io"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/mackerelclient"
)

type channelsApp struct {
	client    mackerelclient.Client
	outStream io.Writer
}

func (app *channelsApp) run() error {
	channels, err := app.client.FindChannels()
	if err != nil {
		return err
	}

	format.PrettyPrintJSON(app.outStream, channels)
	return nil
}

func (app *channelsApp) pullChannels(isVerbose bool, optFilePath string) error {
	// TODO
	return nil
}
