package services

import (
	"io"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type servicesApp struct {
	client    mackerelclient.Client
	outStream io.Writer
	jqFilter  string
}

func (app *servicesApp) run() error {
	services, err := app.client.FindServices()
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, services, app.jqFilter)
	logger.DieIf(err)
	return nil
}
