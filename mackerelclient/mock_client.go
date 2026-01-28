package mackerelclient

import (
	"context"
	"errors"

	"github.com/mackerelio/mackerel-client-go"
)

// MockClient represents a mock client of Mackerel API
type MockClient struct {
	findAWSIntegrationsCallback func() ([]*mackerel.AWSIntegration, error)
	findHostsCallback           func(param *mackerel.FindHostsParam) ([]*mackerel.Host, error)
	findHostCallback            func(id string) (*mackerel.Host, error)
	findServicesCallback        func() ([]*mackerel.Service, error)
	findChannelsCallback        func() ([]*mackerel.Channel, error)
	findUsersCallback           func() ([]*mackerel.User, error)
	getOrgCallback              func() (*mackerel.Org, error)
	createHostCallback          func(param *mackerel.CreateHostParam) (string, error)
	updateHostStatusCallback    func(hostID string, status string) error
	listHostMetricNamesCallback func(id string) ([]string, error)
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
func (c *MockClient) FindHostContext(ctx context.Context, id string) (*mackerel.Host, error) {
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
func (c *MockClient) FindHostsContext(ctx context.Context, param *mackerel.FindHostsParam) ([]*mackerel.Host, error) {
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
func (c *MockClient) FindServicesContext(ctx context.Context) ([]*mackerel.Service, error) {
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
func (c *MockClient) FindChannelsContext(ctx context.Context) ([]*mackerel.Channel, error) {
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
func (c *MockClient) GetOrgContext(ctx context.Context) (*mackerel.Org, error) {
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
func (c *MockClient) CreateHostContext(ctx context.Context, param *mackerel.CreateHostParam) (string, error) {
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
func (c *MockClient) UpdateHostStatusContext(ctx context.Context, hostID string, status string) error {
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
func (c *MockClient) FindAWSIntegrationsContext(ctx context.Context) ([]*mackerel.AWSIntegration, error) {
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
func (c *MockClient) ListHostMetricNamesContext(ctx context.Context, hostID string) ([]string, error) {
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

// FindUsers ...
func (c *MockClient) FindUsersContext(ctx context.Context) ([]*mackerel.User, error) {
	if c.findUsersCallback != nil {
		return c.findUsersCallback()
	}
	return nil, errCallbackNotFound("FindUsers")
}

// MockFindUsers returns an option to set the callback of FindUsers
func MockFindUsers(callback func() ([]*mackerel.User, error)) MockClientOption {
	return func(c *MockClient) {
		c.findUsersCallback = callback
	}
}

var errNotImplemented = errors.New("not implemented")

func (c *MockClient) FindWithClosedAlertsContext(context.Context) (*mackerel.AlertsResp, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FindWithClosedAlertsByNextIDContext(context.Context, string) (*mackerel.AlertsResp, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FindAlertsContext(context.Context) (*mackerel.AlertsResp, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FindAlertsByNextIDContext(context.Context, string) (*mackerel.AlertsResp, error) {
	return nil, errNotImplemented
}
func (c *MockClient) CloseAlertContext(context.Context, string, string) (*mackerel.Alert, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FindAlertLogsContext(context.Context, string, *mackerel.FindAlertLogsParam) (*mackerel.FindAlertLogsResp, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FindMonitorsContext(context.Context) ([]mackerel.Monitor, error) {
	return nil, errNotImplemented
}

func (c *MockClient) CreateGraphAnnotationContext(ctx context.Context, annotation *mackerel.GraphAnnotation) (*mackerel.GraphAnnotation, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FindGraphAnnotationsContext(ctx context.Context, service string, from int64, to int64) ([]*mackerel.GraphAnnotation, error) {
	return nil, errNotImplemented
}
func (c *MockClient) UpdateGraphAnnotationContext(ctx context.Context, annotationID string, annotation *mackerel.GraphAnnotation) (*mackerel.GraphAnnotation, error) {
	return nil, errNotImplemented
}
func (c *MockClient) DeleteGraphAnnotationContext(ctx context.Context, annotationID string) (*mackerel.GraphAnnotation, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FindDashboardsContext(ctx context.Context) ([]*mackerel.Dashboard, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FindDashboardContext(ctx context.Context, dashboardID string) (*mackerel.Dashboard, error) {
	return nil, errNotImplemented
}
func (c *MockClient) UpdateDashboardContext(ctx context.Context, dashboardID string, param *mackerel.Dashboard) (*mackerel.Dashboard, error) {
	return nil, errNotImplemented
}
func (c *MockClient) CreateDashboardContext(ctx context.Context, param *mackerel.Dashboard) (*mackerel.Dashboard, error) {
	return nil, errNotImplemented
}
func (c *MockClient) RetireHostContext(ctx context.Context, hostID string) error {
	return errNotImplemented
}
func (c *MockClient) UpdateHostRoleFullnamesContext(ctx context.Context, hostID string, roleFullnames []string) error {
	return errNotImplemented
}
func (c *MockClient) UpdateHostContext(ctx context.Context, hostID string, param *mackerel.UpdateHostParam) (string, error) {
	return "", errNotImplemented
}
func (c *MockClient) ListServiceMetricNamesContext(ctx context.Context, serviceName string) ([]string, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FetchHostMetricValuesContext(ctx context.Context, hostID string, metricName string, from int64, to int64) ([]mackerel.MetricValue, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FetchServiceMetricValuesContext(ctx context.Context, serviceName string, metricName string, from int64, to int64) ([]mackerel.MetricValue, error) {
	return nil, errNotImplemented
}
func (c *MockClient) FetchLatestMetricValuesContext(ctx context.Context, hostIDs []string, metricNames []string) (mackerel.LatestMetricValues, error) {
	return nil, errNotImplemented
}
func (c *MockClient) PostHostMetricValuesByHostIDContext(ctx context.Context, hostID string, metricValues []*mackerel.MetricValue) error {
	return errNotImplemented
}
func (c *MockClient) PostServiceMetricValuesContext(ctx context.Context, serviceName string, metricValues []*mackerel.MetricValue) error {
	return errNotImplemented
}
func (c *MockClient) CreateMonitorContext(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error) {
	return nil, errNotImplemented
}
func (c *MockClient) DeleteMonitorContext(ctx context.Context, monitorID string) (mackerel.Monitor, error) {
	return nil, errNotImplemented
}
func (c *MockClient) UpdateMonitorContext(ctx context.Context, monitorID string, param mackerel.Monitor) (mackerel.Monitor, error) {
	return nil, errNotImplemented
}
