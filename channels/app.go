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
	jq        string
}

func (app *channelsApp) run() error {
	channels, err := app.client.FindChannels()
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, channels, app.jq)
	logger.DieIf(err)
	return nil
}

func (app *channelsApp) pullChannels(isVerbose bool, optFilePath string) error {
	channels, err := app.client.FindChannels()
	logger.DieIf(err)

	filePath := "channels.json"
	if optFilePath != "" {
		filePath = optFilePath
	}

	err = saveChannels(channels, filePath)
	logger.DieIf(err)

	if isVerbose {
		err := format.PrettyPrintJSON(os.Stdout, channels, "")
		logger.DieIf(err)
	}

	logger.Log("info", fmt.Sprintf("Channels are saved to '%s' (%d rules).", filePath, len(channels)))
	return nil
}

var saveChannels = channelSaveRules

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
