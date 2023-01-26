package metrics_scope_collector

import (
	"context"
	"os"
	"testing"

	metricsscope "cloud.google.com/go/monitoring/metricsscope/apiv1"
	crmbox "github.com/sinmetalcraft/gcpbox/cloudresourcemanager/v3"
	"google.golang.org/api/cloudresourcemanager/v3"
)

func TestService_ImportMonitoredProjects(t *testing.T) {
	ctx := context.Background()

	project := os.Getenv("TEST_GOOGLE_CLOUD_PROJECT_MONITORED_PROJECT") // これscopingProjectだわ

	client, err := metricsscope.NewMetricsScopesClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	metricsScopesService, err := NewMetricsScopesService(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		t.Fatal(err)
	}

	resourceManagerService, err := crmbox.NewResourceManagerService(ctx, crmService)
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewService(ctx, metricsScopesService, resourceManagerService)
	if err != nil {
		t.Fatal(err)
	}

	count, err := s.ImportMonitoredProjects(ctx, project, &crmbox.ResourceID{ID: "190932998497", Type: "organization"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Import MonitoredProject Count %d\n", count)
}
