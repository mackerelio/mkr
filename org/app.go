package org

import (
	"io"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type orgApp struct {
	client    mackerelclient.Client
	outStream io.Writer
}

func (app *orgApp) run() error {
	org, err := app.client.GetOrg()
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, org)
	logger.DieIf(err)
	return nil
}
