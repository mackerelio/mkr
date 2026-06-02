package apm

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/mackerelclient"
)

type appLogger interface {
	Log(string, string)
	Error(error)
}

type httpServerStatsApp struct {
	client    mackerelclient.Client
	outStream io.Writer
	logger    appLogger
}

func (app *httpServerStatsApp) listHTTPServerStats(ctx context.Context, param *mackerel.ListHTTPServerStatsParam) error {
	stats, err := app.client.ListHTTPServerStatsContext(ctx, param)
	if err != nil {
		app.logger.Error(err)
		return err
	}

	encoder := json.NewEncoder(app.outStream)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(stats.Results); err != nil {
		app.logger.Error(err)
		return err
	}
	return nil
}
