package status

import (
	"io"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type statussApp struct {
	client    mackerelclient.Client
	outStream io.Writer

	argHostID string
	isVerbose bool
}

func (app *statussApp) run() error {
	host, err := app.client.FindHost(app.argHostID)
	if err != nil {
		return err
	}

	if app.isVerbose {
		err := format.PrettyPrintJSON(app.outStream, host)
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
		})
		logger.DieIf(err)
	}
	return nil
}
