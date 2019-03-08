package org

import (
	"io"

	"github.com/mackerelio/mkr/format"
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

	format.PrettyPrintJSON(app.outStream, org)
	return nil
}
