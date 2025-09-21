package mackerelclient

import "github.com/mackerelio/mackerel-client-go"

// MockClient represents a mock client of Mackerel API
type MockClient struct {
	findAWSIntegrationsCallback func() ([]*mackerel.AWSIntegration, error)
	findHostsCallback           func(param *mackerel.FindHostsParam) ([]*mackerel.Host, error)
	findHostCallback            func(id string) (*mackerel.Host, error)
	findServicesCallback        func() ([]*mackerel.Service, error)
	findChannelsCallback        func() ([]*mackerel.Channel, error)
	getOrgCallback              func() (*mackerel.Org, error)
	createHostCallback          func(param *mackerel.CreateHostParam) (string, error)
	updateHostStatusCallback    func(hostID string, status string) error
	listHostMetricNamesCallback func(id string) ([]string, error)
	getTraceCallback            func(traceID string) (*mackerel.TraceResponse, error)
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

// FindHost ...
func (c *MockClient) FindHost(id string) (*mackerel.Host, error) {
	if c.findHostCallback != nil {
		return c.findHostCallback(id)
	}
	return nil, errCallbackNotFound("FindHost")
}

// MockFindHost returns an option to set the callback of FindHost
func MockFindHost(callback func(id string) (*mackerel.Host, error)) MockClientOption {
	return func(c *MockClient) {
		c.findHostCallback = callback
	}
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

// FindServices ...
func (c *MockClient) FindServices() ([]*mackerel.Service, error) {
	if c.findServicesCallback != nil {
		return c.findServicesCallback()
	}
	return nil, errCallbackNotFound("FindServices")
}

// MockFindServices returns an option to set the callback of FindServices
func MockFindServices(callback func() ([]*mackerel.Service, error)) MockClientOption {
	return func(c *MockClient) {
		c.findServicesCallback = callback
	}
}

// FindChannels ...
func (c *MockClient) FindChannels() ([]*mackerel.Channel, error) {
	if c.findChannelsCallback != nil {
		return c.findChannelsCallback()
	}
	return nil, errCallbackNotFound("FindChannels")
}

// MockFindChannels returns an option to set the callback of FindChannels
func MockFindChannels(callback func() ([]*mackerel.Channel, error)) MockClientOption {
	return func(c *MockClient) {
		c.findChannelsCallback = callback
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

// CreateHost ...
func (c *MockClient) CreateHost(param *mackerel.CreateHostParam) (string, error) {
	if c.createHostCallback != nil {
		return c.createHostCallback(param)
	}
	return "", errCallbackNotFound("CreateHost")
}

// MockCreateHost returns an option to set the callback of CreateHost
func MockCreateHost(callback func(*mackerel.CreateHostParam) (string, error)) MockClientOption {
	return func(c *MockClient) {
		c.createHostCallback = callback
	}
}

// UpdateHostStatus ...
func (c *MockClient) UpdateHostStatus(hostID string, status string) error {
	if c.updateHostStatusCallback != nil {
		return c.updateHostStatusCallback(hostID, status)
	}
	return errCallbackNotFound("UpdateHostStatus")
}

// MockUpdateHostStatus returns an option to set the callback of UpdateHostStatus
func MockUpdateHostStatus(callback func(string, string) error) MockClientOption {
	return func(c *MockClient) {
		c.updateHostStatusCallback = callback
	}
}

// FindAWSIntegrations ...
func (c *MockClient) FindAWSIntegrations() ([]*mackerel.AWSIntegration, error) {
	if c.findAWSIntegrationsCallback != nil {
		return c.findAWSIntegrationsCallback()
	}
	return nil, errCallbackNotFound("FindAWSIntegrations")
}

// MockFindAWSIntegrations returns an option to set the callback of FindAWSIntegrations
func MockFindAWSIntegrations(callback func() ([]*mackerel.AWSIntegration, error)) MockClientOption {
	return func(c *MockClient) {
		c.findAWSIntegrationsCallback = callback
	}
}

// ListHostMetricNames ...
func (c *MockClient) ListHostMetricNames(hostID string) ([]string, error) {
	if c.listHostMetricNamesCallback != nil {
		return c.listHostMetricNamesCallback(hostID)
	}
	return nil, errCallbackNotFound("ListHostMetricNames")
}

// MockListHostMetricNames returns an option to set the callback of ListHostMetricNames
func MockListHostMetricNames(callback func(string) ([]string, error)) MockClientOption {
	return func(c *MockClient) {
		c.listHostMetricNamesCallback = callback
	}
}

// GetTrace ...
func (c *MockClient) GetTrace(traceID string) (*mackerel.TraceResponse, error) {
	if c.getTraceCallback != nil {
		return c.getTraceCallback(traceID)
	}
	return nil, errCallbackNotFound("GetTrace")
}

// MockGetTrace returns an option to set the callback of GetTrace
func MockGetTrace(callback func(traceID string) (*mackerel.TraceResponse, error)) MockClientOption {
	return func(c *MockClient) {
		c.getTraceCallback = callback
	}
}
