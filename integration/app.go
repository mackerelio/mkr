package integration

import (
	"io"
	"strings"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type integrationApp struct {
	client    mackerelclient.Client
	outStream io.Writer
}

type integrations struct {
	AWS []*mackerel.AWSIntegration
}

func (app *integrationApp) run(provider string) error {
	var providers []string
	if provider == "" || provider == "all" {
		providers = getMackerelAPISupportProvider()
	} else {
		providers = strings.Split(provider, ",")
	}

	integration := integrations{
		AWS: []*mackerel.AWSIntegration{},
	}
	var err error
	for _, p := range providers {
		switch p {
		case "aws":
			integration.AWS, err = app.findAWSIntegrations()
			if err != nil {
				return err
			}
		}
	}

	err = format.PrettyPrintJSON(app.outStream, integration)
	logger.DieIf(err)
	return nil
}

func getMackerelAPISupportProvider() []string {
	return []string{"aws"}
}

func (app *integrationApp) findAWSIntegrations() ([]*mackerel.AWSIntegration, error) {
	awsIntegrations, err := app.client.FindAWSIntegrations()
	if err != nil {
		return nil, err
	}
	return awsIntegrations, nil
}
