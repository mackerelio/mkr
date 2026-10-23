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

func (app *logsApp) search(ctx context.Context, param *mackerel.ListLogsParam) error {
	res, err := app.client.ListLogsContext(ctx, param)
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, res, app.jqFilter)
	logger.DieIf(err)
	return nil
}

func (app *logsApp) listSavedSearches(ctx context.Context) error {
	res, err := app.client.FindSavedLogSearchesContext(ctx)
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, res, app.jqFilter)
	logger.DieIf(err)
	return nil
}
