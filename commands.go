package main

import (
	"github.com/mackerelio/mkr/alerts"
	"github.com/mackerelio/mkr/annotations"
	"github.com/mackerelio/mkr/aws_integrations"
	"github.com/mackerelio/mkr/channels"
	"github.com/mackerelio/mkr/checks"
	"github.com/mackerelio/mkr/dashboards"
	"github.com/mackerelio/mkr/hosts"
	"github.com/mackerelio/mkr/metric_names"
	"github.com/mackerelio/mkr/metrics"
	"github.com/mackerelio/mkr/monitors"
	"github.com/mackerelio/mkr/org"
	"github.com/mackerelio/mkr/plugin"
	"github.com/mackerelio/mkr/services"
	"github.com/mackerelio/mkr/status"
	"github.com/mackerelio/mkr/wrap"
	"github.com/urfave/cli"
)

// Commands cli.Command object list
var Commands = []cli.Command{
	status.Command,
	hosts.CommandHosts,
	hosts.CommandCreate,
	hosts.CommandUpdate,
	metrics.CommandThrow,
	metrics.Command,
	metrics.CommandFetch,
	hosts.CommandRetire,
	services.Command,
	monitors.Command,
	channels.Command,
	alerts.Command,
	dashboards.Command,
	annotations.Command,
	org.Command,
	plugin.CommandPlugin,
	checks.Command,
	wrap.Command,
	aws_integrations.Command,
	metric_names.Command,
}
