package main

import (
	"github.com/mackerelio/mkr/aws_integrations"
	"github.com/mackerelio/mkr/channels"
	"github.com/mackerelio/mkr/checks"
	"github.com/mackerelio/mkr/fetch"
	"github.com/mackerelio/mkr/hosts"
	"github.com/mackerelio/mkr/metrics"
	"github.com/mackerelio/mkr/monitors"
	"github.com/mackerelio/mkr/org"
	"github.com/mackerelio/mkr/plugin"
	"github.com/mackerelio/mkr/retire"
	"github.com/mackerelio/mkr/services"
	"github.com/mackerelio/mkr/status"
	"github.com/mackerelio/mkr/throw"
	"github.com/mackerelio/mkr/update"
	"github.com/mackerelio/mkr/wrap"
	"github.com/urfave/cli"
)

// Commands cli.Command object list
var Commands = []cli.Command{
	status.Command,
	hosts.CommandHosts,
	hosts.CommandCreate,
	update.Command,
	throw.Command,
	metrics.Command,
	fetch.Command,
	retire.Command,
	services.Command,
	monitors.Command,
	channels.Command,
	commandAlerts,
	commandDashboards,
	commandAnnotations,
	org.Command,
	plugin.CommandPlugin,
	checks.Command,
	wrap.Command,
	aws_integrations.Command,
}
