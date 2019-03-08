package mackerelclient

import mackerel "github.com/mackerelio/mackerel-client-go"

// MockClient represents a mock client of Mackerel API
type MockClient struct {
	findHostsCallback func(param *mackerel.FindHostsParam) ([]*mackerel.Host, error)
	getOrgCallback    func() (*mackerel.Org, error)
}

// MockClientOption represents an option of mock client of Mackerel API
type MockClientOption func(*MockClient)

// NewMockClient creates a new mock client of Mackerel API
func NewMockClient(opts ...MockClientOption) *MockClient {
	client := &MockClient{}
	for _, opt := range opts {
		client.ApplyOption(opt)
	}
	return client
}

// ApplyOption applies a mock client option
func (c *MockClient) ApplyOption(opt MockClientOption) {
	opt(c)
}

type errCallbackNotFound string

func (err errCallbackNotFound) Error() string {
	return string(err) + " callback not found"
}

// FindHosts ...
func (c *MockClient) FindHosts(param *mackerel.FindHostsParam) ([]*mackerel.Host, error) {
	if c.findHostsCallback != nil {
		return c.findHostsCallback(param)
	}
	return nil, errCallbackNotFound("FindHosts")
}

// MockFindHosts returns an option to set the callback of FindHosts
func MockFindHosts(callback func(param *mackerel.FindHostsParam) ([]*mackerel.Host, error)) MockClientOption {
	return func(c *MockClient) {
		c.findHostsCallback = callback
	}
}

// GetOrg ...
func (c *MockClient) GetOrg() (*mackerel.Org, error) {
	if c.getOrgCallback != nil {
		return c.getOrgCallback()
	}
	return nil, errCallbackNotFound("GetOrg")
}

// MockGetOrg returns an option to set the callback of GetOrg
func MockGetOrg(callback func() (*mackerel.Org, error)) MockClientOption {
	return func(c *MockClient) {
		c.getOrgCallback = callback
	}
}
