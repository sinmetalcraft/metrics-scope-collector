package metrics_scope_collector

import (
	"context"
	"fmt"
	"strings"

	crmbox "github.com/sinmetalcraft/gcpbox/cloudresourcemanager/v3"
)

type Service struct {
	MetricsScopesService   *MetricsScopesService
	ResourceManagerService *crmbox.ResourceManagerService
}

func NewService(ctx context.Context, metricsScopesService *MetricsScopesService, resourceManagerService *crmbox.ResourceManagerService) (*Service, error) {
	return &Service{
		MetricsScopesService:   metricsScopesService,
		ResourceManagerService: resourceManagerService,
	}, nil
}

// ImportMonitoredProjects is scopingProjectのMetricsScopeにparentResourceID配下のProjectを追加する
func (s *Service) ImportMonitoredProjects(ctx context.Context, scopingProject string, parentResourceID *crmbox.ResourceID) (int, error) {
	scope, err := s.MetricsScopesService.GetMetricsScope(ctx, scopingProject)
	if err != nil {
		return 0, fmt.Errorf("failed MetricsScopesService.GetMetricsScope. scopingProject=%s,parentResourceID=%v : %w", scopingProject, parentResourceID, err)
	}
	existsMonitoredProjects := make(map[string]bool, len(scope.GetMonitoredProjects()))
	for _, v := range scope.GetMonitoredProjects() {
		// locations/global/metricsScopes/{ScopingProjectNumber}/projects/{MonitoredProjectNumber}
		l := strings.Split(v.GetName(), "/")
		if len(l) != 6 {
			return 0, fmt.Errorf("invalid MonitoredProjects format. %s", v.GetName())
		}
		existsMonitoredProjects[l[5]] = true
	}

	l, err := s.ResourceManagerService.GetRelatedProject(ctx, parentResourceID)
	if err != nil {
		return 0, fmt.Errorf("failed ResourceManagerService.GetRelatedProject. scopingProject=%s,parentResourceID=%v : %w", scopingProject, parentResourceID, err)
	}

	var createdCount int
	for _, v := range l {
		projectNumber := strings.ReplaceAll(v.Name, "projects/", "")
		if projectNumber == scopingProject {
			continue
		}
		_, ok := existsMonitoredProjects[projectNumber]
		if ok {
			continue
		}

		ret, err := s.MetricsScopesService.CreateMonitoredProject(ctx, scopingProject, projectNumber)
		if err != nil {
			fmt.Printf("failed CreateMonitoredProject: %s(%s). %s\n", v.ProjectId, projectNumber, err)
			continue
		}
		fmt.Printf("created MonitoredProject: %s (%s)\n", ret.GetName(), v.ProjectId)
		createdCount++
	}
	return createdCount, nil
}

// CleanUp is 指定したScopingProjectのMetrics Scopeをすべて削除して最初の状態に戻す
func (s *Service) CleanUp(ctx context.Context, scopingProject string) error {
	scope, err := s.MetricsScopesService.GetMetricsScope(ctx, scopingProject)
	if err != nil {
		return fmt.Errorf("failed MetricsScopesService.GetMetricsScope. scopingProject=%s : %w", scopingProject, err)
	}

	for _, v := range scope.GetMonitoredProjects() {
		if err := s.MetricsScopesService.DeleteMonitoredProject(ctx, scopingProject, v.GetName()); err != nil {
			return err
		}
	}
	return nil
}
