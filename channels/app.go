package channels

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
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
	channels, err := mackerelclient.NewFromContext(c).FindChannels()
	logger.DieIf(err)

	filePath := "channels.json"
	if optFilePath != "" {
		filePath = optFilePath
	}

	channelSaveRules(channels, filePath)

	if isVerbose {
		format.PrettyPrintJSON(os.Stdout, channels)
	}

	logger.Log("info", fmt.Sprintf("Channels are saved to '%s' (%d rules).", filePath, len(channels)))
	return nil
}

func channelSaveRules(rules []*mackerel.Channel, filePath string) error {
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
