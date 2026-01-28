package status

import (
	"context"
	"io"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type statussApp struct {
	client    mackerelclient.Client
	outStream io.Writer

	argHostID string
	isVerbose bool
	jqFilter  string
}

type HostWithMetrics struct {
	*mackerel.Host
	Metrics []string `json:"metrics,omitempty"`
}

func (app *statussApp) run(ctx context.Context) error {
	host, err := app.client.FindHostContext(ctx, app.argHostID)
	if err != nil {
		return err
	}

	if app.isVerbose {
		metrics, err := app.client.ListHostMetricNamesContext(ctx, host.ID)
		logger.DieIf(err)
		hostWithMetrics := HostWithMetrics{Host: host, Metrics: metrics}
		err = format.PrettyPrintJSON(app.outStream, hostWithMetrics, app.jqFilter)
		logger.DieIf(err)
	} else {
		err := format.PrettyPrintJSON(app.outStream, &format.Host{
			ID:            host.ID,
			Name:          host.Name,
			DisplayName:   host.DisplayName,
			Status:        host.Status,
			RoleFullnames: host.GetRoleFullnames(),
			IsRetired:     host.IsRetired,
			CreatedAt:     format.ISO8601Extended(host.DateFromCreatedAt()),
			IPAddresses:   host.IPAddresses(),
		}, app.jqFilter)
		logger.DieIf(err)
	}
	return nil
}
