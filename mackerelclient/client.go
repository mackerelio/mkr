package mackerelclient

import mackerel "github.com/mackerelio/mackerel-client-go"

// Client represents a client of Mackerel API
type Client interface {
	FindHosts(param *mackerel.FindHostsParam) ([]*mackerel.Host, error)
	FindServices() ([]*mackerel.Service, error)
	GetOrg() (*mackerel.Org, error)
	CreateHost(param *mackerel.CreateHostParam) (string, error)
	UpdateHostStatus(hostID string, status string) error
}
