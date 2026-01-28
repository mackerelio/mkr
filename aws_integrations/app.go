package aws_integrations

import (
	"context"
	"io"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type awsIntegrationsApp struct {
	client    mackerelclient.Client
	outStream io.Writer
	jqFilter  string
}

func (app *awsIntegrationsApp) run(ctx context.Context) error {
	awsIntegrations, err := app.client.FindAWSIntegrationsContext(ctx)
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, awsIntegrations, app.jqFilter)
	logger.DieIf(err)
	return nil
}
