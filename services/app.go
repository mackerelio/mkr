package services

import (
	"io"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/mackerelclient"
)

type servicesApp struct {
	client    mackerelclient.Client
	outStream io.Writer
}

func (app *servicesApp) run() error {
	services, err := app.client.FindServices()
	if err != nil {
		return err
	}

	format.PrettyPrintJSON(app.outStream, services)
	return nil
}
