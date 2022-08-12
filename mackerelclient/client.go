package mackerelclient

import "github.com/mackerelio/mackerel-client-go"

// Client represents a client of Mackerel API
type Client interface {
	FindAWSIntegrations() ([]*mackerel.AWSIntegration, error)
	FindHosts(param *mackerel.FindHostsParam) ([]*mackerel.Host, error)
	FindHost(id string) (*mackerel.Host, error)
	FindServices() ([]*mackerel.Service, error)
	FindChannels() ([]*mackerel.Channel, error)
	GetOrg() (*mackerel.Org, error)
	CreateHost(param *mackerel.CreateHostParam) (string, error)
	UpdateHostStatus(hostID string, status string) error
}
