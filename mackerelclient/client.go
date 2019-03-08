package mackerelclient

import mackerel "github.com/mackerelio/mackerel-client-go"

// Client represents a client of Mackerel API
type Client interface {
	FindHosts(param *mackerel.FindHostsParam) ([]*mackerel.Host, error)
}
