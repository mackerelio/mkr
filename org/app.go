package org

import (
	"context"
	"io"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type orgApp struct {
	client    mackerelclient.Client
	outStream io.Writer
	jqFilter  string
}

func (app *orgApp) run(ctx context.Context) error {
	org, err := app.client.GetOrgContext(ctx)
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, org, app.jqFilter)
	logger.DieIf(err)
	return nil
}
