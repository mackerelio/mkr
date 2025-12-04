package channels

import (
	"encoding/json"
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
	jqFilter  string
}

const defaultFilePath = "channels.json"

func (app *channelsApp) run() error {
	channels, err := app.client.FindChannels()
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, channels, app.jqFilter)
	logger.DieIf(err)
	return nil
}

func (app *channelsApp) pullChannels(isVerbose bool, optFilePath string) error {
	channels, err := app.client.FindChannels()
	logger.DieIf(err)

	filePath := defaultFilePath
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

func (app *channelsApp) pushChannels(isVerbose bool, optFilePath string) error {
	filePath := defaultFilePath
	if optFilePath != "" {
		filePath = optFilePath
	}

	localChannels, err := channelLoadChannels(filePath)
	logger.DieIf(err)

	remoteChannels, err := app.client.FindChannels()
	logger.DieIf(err)

	remoteChannelMap := make(map[string]*mackerel.Channel)
	for _, rc := range remoteChannels {
		if rc.ID != "" {
			remoteChannelMap[rc.ID] = rc
		}
	}

	for _, lc := range localChannels {
		if lc.ID != "" {
			if rc, ok := remoteChannelMap[lc.ID]; ok {
				_, err := app.client.UpdateChannel(rc.ID, lc)
				logger.DieIf(err)
				continue
			}
			logger.Log("info", fmt.Sprintf("Channel ID '%s' not found. Creating a new channel.", lc.ID))
		}
		_, err := app.client.CreateChannel(lc)
		logger.DieIf(err)
	}

	if isVerbose {
		err := format.PrettyPrintJSON(os.Stdout, localChannels, "")
		logger.DieIf(err)
	}
	return nil
}

func channelLoadChannels(filePath string) ([]*mackerel.Channel, error) {
	src, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	channelsWrapper := struct {
		Channels []*mackerel.Channel `json:"channels"`
	}{}

	err = json.NewDecoder(src).Decode(&channelsWrapper)
	if err != nil {
		return nil, err
	}
	return channelsWrapper.Channels, nil
}
