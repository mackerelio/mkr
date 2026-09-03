package logs

import (
	"context"
	"io"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type logsApp struct {
	client    mackerelclient.Client
	outStream io.Writer
	jqFilter  string
}

func (app *logsApp) findLogs(ctx context.Context, param *mackerel.FindLogsParam) error {
	resp, err := app.client.FindLogsContext(ctx, param)
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, resp, app.jqFilter)
	logger.DieIf(err)
	return nil
}
