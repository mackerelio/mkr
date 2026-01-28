package mackerelclient

import (
	"context"

	"github.com/mackerelio/mackerel-client-go"
)

// Client represents a client of Mackerel API
type Client interface {
	FindAWSIntegrationsContext(ctx context.Context) ([]*mackerel.AWSIntegration, error)
	FindHostsContext(ctx context.Context, param *mackerel.FindHostsParam) ([]*mackerel.Host, error)
	FindHostContext(ctx context.Context, id string) (*mackerel.Host, error)
	FindServicesContext(ctx context.Context) ([]*mackerel.Service, error)
	FindChannelsContext(ctx context.Context) ([]*mackerel.Channel, error)
	FindUsersContext(ctx context.Context) ([]*mackerel.User, error)
	GetOrgContext(ctx context.Context) (*mackerel.Org, error)
	CreateHostContext(ctx context.Context, param *mackerel.CreateHostParam) (string, error)
	UpdateHostStatusContext(ctx context.Context, hostID string, status string) error
	ListHostMetricNamesContext(ctx context.Context, id string) ([]string, error)
	// below mock needed implemented.
	FindWithClosedAlertsContext(ctx context.Context) (*mackerel.AlertsResp, error)
	FindWithClosedAlertsByNextIDContext(ctx context.Context, nextID string) (*mackerel.AlertsResp, error)
	FindAlertsContext(ctx context.Context) (*mackerel.AlertsResp, error)
	FindAlertsByNextIDContext(ctx context.Context, nextID string) (*mackerel.AlertsResp, error)
	CloseAlertContext(ctx context.Context, alertID string, reason string) (*mackerel.Alert, error)
	FindAlertLogsContext(ctx context.Context, alertId string, params *mackerel.FindAlertLogsParam) (*mackerel.FindAlertLogsResp, error)
	FindMonitorsContext(ctx context.Context) ([]mackerel.Monitor, error)
	CreateGraphAnnotationContext(ctx context.Context, annotation *mackerel.GraphAnnotation) (*mackerel.GraphAnnotation, error)
	FindGraphAnnotationsContext(ctx context.Context, service string, from int64, to int64) ([]*mackerel.GraphAnnotation, error)
	UpdateGraphAnnotationContext(ctx context.Context, annotationID string, annotation *mackerel.GraphAnnotation) (*mackerel.GraphAnnotation, error)
	DeleteGraphAnnotationContext(ctx context.Context, annotationID string) (*mackerel.GraphAnnotation, error)
	FindDashboardsContext(ctx context.Context) ([]*mackerel.Dashboard, error)
	FindDashboardContext(ctx context.Context, dashboardID string) (*mackerel.Dashboard, error)
	UpdateDashboardContext(ctx context.Context, dashboardID string, param *mackerel.Dashboard) (*mackerel.Dashboard, error)
	CreateDashboardContext(ctx context.Context, param *mackerel.Dashboard) (*mackerel.Dashboard, error)
	RetireHostContext(ctx context.Context, hostID string) error
	UpdateHostRoleFullnamesContext(ctx context.Context, hostID string, roleFullnames []string) error
	UpdateHostContext(ctx context.Context, hostID string, param *mackerel.UpdateHostParam) (string, error)
	ListServiceMetricNamesContext(ctx context.Context, serviceName string) ([]string, error)
	FetchHostMetricValuesContext(ctx context.Context, hostID string, metricName string, from int64, to int64) ([]mackerel.MetricValue, error)
	FetchServiceMetricValuesContext(ctx context.Context, serviceName string, metricName string, from int64, to int64) ([]mackerel.MetricValue, error)
	FetchLatestMetricValuesContext(ctx context.Context, hostIDs []string, metricNames []string) (mackerel.LatestMetricValues, error)
	PostHostMetricValuesByHostIDContext(ctx context.Context, hostID string, metricValues []*mackerel.MetricValue) error
	PostServiceMetricValuesContext(ctx context.Context, serviceName string, metricValues []*mackerel.MetricValue) error
	CreateMonitorContext(ctx context.Context, param mackerel.Monitor) (mackerel.Monitor, error)
	DeleteMonitorContext(ctx context.Context, monitorID string) (mackerel.Monitor, error)
	UpdateMonitorContext(ctx context.Context, monitorID string, param mackerel.Monitor) (mackerel.Monitor, error)
}
