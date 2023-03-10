package metrics_scope_collector

import (
	"context"
	"fmt"

	metricsscope "cloud.google.com/go/monitoring/metricsscope/apiv1"
	metricsscopepb "cloud.google.com/go/monitoring/metricsscope/apiv1/metricsscopepb"
)

type MetricsScopesService struct {
	metricsScopeClient *metricsscope.MetricsScopesClient
}

func NewMetricsScopesService(ctx context.Context, metricsScopeClient *metricsscope.MetricsScopesClient) (*MetricsScopesService, error) {
	return &MetricsScopesService{
		metricsScopeClient,
	}, nil
}

// ListMetricsScopesByMonitoredProject is 指定したProjectのMetricsScopesを返す
// 指定するのはPROJECT_ID or PROJECT_NUMBER
func (s *MetricsScopesService) ListMetricsScopesByMonitoredProject(ctx context.Context, project string) ([]*metricsscopepb.MetricsScope, error) {
	req := &metricsscopepb.ListMetricsScopesByMonitoredProjectRequest{
		MonitoredResourceContainer: fmt.Sprintf("projects/%s", project),
	}
	resp, err := s.metricsScopeClient.ListMetricsScopesByMonitoredProject(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetMetricsScopes(), nil
}

// GetMetricsScope is 指定したScopingProjectのMetricsScopeを返す
// 指定するのはPROJECT_ID or PROJECT_NUMBER
func (s *MetricsScopesService) GetMetricsScope(ctx context.Context, project string) (*metricsscopepb.MetricsScope, error) {
	req := &metricsscopepb.GetMetricsScopeRequest{
		Name: fmt.Sprintf("locations/global/metricsScopes/%s", project),
	}
	v, err := s.metricsScopeClient.GetMetricsScope(ctx, req)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// CreateMonitoredProject is scopingProjectにmonitoringProjectのmetricsを追加する
// scopingProject, monitoringProjectはPROJECT_ID or PROJECT_NUMBERを指定する
func (s *MetricsScopesService) CreateMonitoredProject(ctx context.Context, scopingProject string, monitoredProject string) (*metricsscopepb.MonitoredProject, error) {
	req := &metricsscopepb.CreateMonitoredProjectRequest{
		Parent: fmt.Sprintf("locations/global/metricsScopes/%s", scopingProject),
		MonitoredProject: &metricsscopepb.MonitoredProject{
			Name: fmt.Sprintf("locations/global/metricsScopes/%s/projects/%s", scopingProject, monitoredProject),
		},
	}
	ope, err := s.metricsScopeClient.CreateMonitoredProject(ctx, req)
	if err != nil {
		return nil, err
	}
	ret, err := ope.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// DeleteMonitoredProject is 指定したMonitoredProjectをScoping Projectのmetrics scopeから削除する
// scopingProject, monitoringProjectはPROJECT_ID or PROJECT_NUMBERを指定する
func (s *MetricsScopesService) DeleteMonitoredProject(ctx context.Context, scopingProject string, monitoredProject string) error {
	req := &metricsscopepb.DeleteMonitoredProjectRequest{
		Name: fmt.Sprintf("locations/global/metricsScopes/%s/projects/%s", scopingProject, monitoredProject),
	}
	ope, err := s.metricsScopeClient.DeleteMonitoredProject(ctx, req)
	if err != nil {
		return err
	}
	err = ope.Wait(ctx)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMonitoredProjectByMonitoredProjectName is 指定したMonitoredProjectを削除する
//
//	Example:
//	  `locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}/projects/{MONITORED_PROJECT_ID_OR_NUMBER}`
func (s *MetricsScopesService) DeleteMonitoredProjectByMonitoredProjectName(ctx context.Context, monitoredProjectName string) error {
	req := &metricsscopepb.DeleteMonitoredProjectRequest{
		Name: monitoredProjectName,
	}
	ope, err := s.metricsScopeClient.DeleteMonitoredProject(ctx, req)
	if err != nil {
		return err
	}
	err = ope.Wait(ctx)
	if err != nil {
		return err
	}
	return nil
}
