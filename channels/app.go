package channels

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
