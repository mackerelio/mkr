package hosts

import (
	"fmt"
	"io"
	"text/template"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/mackerelclient"
)

type appLogger interface {
	Log(string, string)
	Error(error)
}

type hostApp struct {
	client    mackerelclient.Client
	logger    appLogger
	outStream io.Writer
}

type findHostsParam struct {
	verbose bool

	name     string
	service  string
	roles    []string
	statuses []string

	format string
}

func (ha *hostApp) findHosts(param findHostsParam) error {
	hosts, err := ha.client.FindHosts(&mackerel.FindHostsParam{
		Name:     param.name,
		Service:  param.service,
		Roles:    param.roles,
		Statuses: param.statuses,
	})
	if err != nil {
		return err
	}

	switch {
	case param.format != "":
		t, err := template.New("format").Parse(param.format)
		if err != nil {
			return err
		}
		return t.Execute(ha.outStream, hosts)
	case param.verbose:
		return format.PrettyPrintJSON(ha.outStream, hosts)
	default:
		hostsFormat := make([]*format.Host, 0)
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

type createHostParam struct {
	name             string
	roleFullnames    []string
	status           string
	customIdentifier string
}

func (ha *hostApp) createHost(param createHostParam) error {
	hostID, err := ha.client.CreateHost(&mackerel.CreateHostParam{
		Name:             param.name,
		RoleFullnames:    param.roleFullnames,
		CustomIdentifier: param.customIdentifier,
	})
	if err != nil {
		ha.error(err)
		return err
	}

	ha.log("created", hostID)

	if param.status != "" {
		err := ha.client.UpdateHostStatus(hostID, param.status)
		if err != nil {
			ha.error(err)
			return err
		}
		ha.log("updated", fmt.Sprintf("%s %s", hostID, param.status))
	}
	return nil
}

func (ha *hostApp) log(prefix, message string) {
	ha.logger.Log(prefix, message)
}

func (ha *hostApp) error(err error) {
	ha.logger.Error(err)
}
