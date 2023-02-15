package metrics_scope_collector

import (
	"context"
	"fmt"
	"strings"

	crmbox "github.com/sinmetalcraft/gcpbox/cloudresourcemanager/v3"
	metricsscopebox "github.com/sinmetalcraft/gcpbox/monitoring/metricsscope/v0"
)

type Service struct {
	MetricsScopeService    *metricsscopebox.Service
	ResourceManagerService *crmbox.ResourceManagerService
}

func NewService(ctx context.Context, metricsScopeService *metricsscopebox.Service, resourceManagerService *crmbox.ResourceManagerService) (*Service, error) {
	return &Service{
		MetricsScopeService:    metricsScopeService,
		ResourceManagerService: resourceManagerService,
	}, nil
}

// ImportMonitoredProjects is scopingProjectのMetricsScopeにparentResourceID配下のProjectを追加する
func (s *Service) ImportMonitoredProjects(ctx context.Context, scopingProject string, parentResourceID *crmbox.ResourceID) (int, error) {
	scope, err := s.MetricsScopeService.GetMetricsScope(ctx, scopingProject)
	if err != nil {
		return 0, fmt.Errorf("failed MetricsScopesService.GetMetricsScope. scopingProject=%s,parentResourceID=%v : %w", scopingProject, parentResourceID, err)
	}
	existsMonitoredProjects := make(map[string]bool, len(scope.MonitoredProjects))
	for _, v := range scope.MonitoredProjects {
		// locations/global/metricsScopes/{ScopingProjectNumber}/projects/{MonitoredProjectNumber}
		l := strings.Split(v.Name, "/")
		if len(l) != 6 {
			return 0, fmt.Errorf("invalid MonitoredProjects format. %s", v.Name)
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

		ret, err := s.MetricsScopeService.CreateMonitoredProject(ctx, scopingProject, projectNumber)
		if err != nil {
			fmt.Printf("failed CreateMonitoredProject: %s(%s). %s\n", v.ProjectId, projectNumber, err)
			continue
		}
		fmt.Printf("created MonitoredProject: %s (%s)\n", ret.Name, v.ProjectId)
		createdCount++
	}
	return createdCount, nil
}

// CleanUp is 指定したScopingProjectのMetrics Scopeをすべて削除して最初の状態に戻す
func (s *Service) CleanUp(ctx context.Context, scopingProject string) (int, error) {
	scope, err := s.MetricsScopeService.GetMetricsScope(ctx, scopingProject)
	if err != nil {
		return 0, fmt.Errorf("failed MetricsScopesService.GetMetricsScope. scopingProject=%s : %w", scopingProject, err)
	}

	var count int
	for _, v := range scope.MonitoredProjects {
		scopingProject, err := v.ScopingProjectIDOrNumber()
		if err != nil {
			err = fmt.Errorf("failed ScopingProjectIDOrNumber name=%s : %w", v.Name, err)
			fmt.Printf("%s\n", err)
			continue
		}

		monitoredProject, err := v.MonitoredProjectIDOrNumber()
		if err != nil {
			err = fmt.Errorf("failed MonitoredProjectIDOrNumber name=%s : %w", v.Name, err)
			fmt.Printf("%s\n", err)
			continue
		}

		if scopingProject == monitoredProject {
			// 自分は削除できないので、skip
			continue
		}

		if err := s.MetricsScopeService.DeleteMonitoredProjectByMonitoredProjectName(ctx, v.Name); err != nil {
			err = fmt.Errorf("failed GathererService.DeleteMonitoredProjectByMonitoredProjectName. name=%s : %w", v.Name, err)
			fmt.Printf("%s\n", err)
		} else {
			count++
		}
	}
	return count, err
}
