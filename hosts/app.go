package hosts

import (
	"io"
	"text/template"

	mackerel "github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/mackerelclient"
)

type hostApp struct {
	cli mackerelclient.Client

	verbose bool

	name     string
	service  string
	roles    []string
	statuses []string

	format string

	outStream io.Writer
}

func (ha *hostApp) run() error {
	hosts, err := ha.cli.FindHosts(&mackerel.FindHostsParam{
		Name:     ha.name,
		Service:  ha.service,
		Roles:    ha.roles,
		Statuses: ha.statuses,
	})
	if err != nil {
		return err
	}

	switch {
	case ha.format != "":
		t, err := template.New("format").Parse(ha.format)
		if err != nil {
			return err
		}
		return t.Execute(ha.outStream, hosts)
	case ha.verbose:
		return format.PrettyPrintJSON(ha.outStream, hosts)
	default:
		var hostsFormat []*format.Host
		for _, host := range hosts {
			hostsFormat = append(hostsFormat, &format.Host{
				ID:            host.ID,
				Name:          host.Name,
				DisplayName:   host.DisplayName,
				Status:        host.Status,
				RoleFullnames: host.GetRoleFullnames(),
				IsRetired:     host.IsRetired,
				CreatedAt:     format.ISO8601Extended(host.DateFromCreatedAt()),
				IPAddresses:   host.IPAddresses(),
			})
		}
		return format.PrettyPrintJSON(ha.outStream, hostsFormat)
	}
}
