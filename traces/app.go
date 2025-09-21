package traces

import (
	"io"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
)

type tracesApp struct {
	client    mackerelclient.Client
	outStream io.Writer
	jqFilter  string
}

func (app *tracesApp) getTrace(traceID string) error {
	trace, err := app.client.GetTrace(traceID)
	if err != nil {
		return err
	}

	err = format.PrettyPrintJSON(app.outStream, trace, app.jqFilter)
	logger.DieIf(err)
	return nil
}
