package metrics_scope_collector

import (
	"context"
	"os"
	"testing"

	metricsscope "cloud.google.com/go/monitoring/metricsscope/apiv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMetricsScopesService_ListMetricsScopesByMonitoredProject(t *testing.T) {
	ctx := context.Background()

	project := os.Getenv("TEST_GOOGLE_CLOUD_PROJECT_MONITORED_PROJECT") // これscopingProjectだわ

	client, err := metricsscope.NewMetricsScopesClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewMetricsScopesService(ctx, client)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.ListMetricsScopesByMonitoredProject(ctx, project)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range got {
		t.Logf("%s\n", v.GetName())
		for _, p := range v.GetMonitoredProjects() {
			t.Logf("\t%s\n", p.GetName())
		}
	}
}

func TestMetricsScopesService_GetMetricsScope(t *testing.T) {
	ctx := context.Background()

	project := os.Getenv("TEST_GOOGLE_CLOUD_PROJECT_MONITORED_PROJECT")

	client, err := metricsscope.NewMetricsScopesClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewMetricsScopesService(ctx, client)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.GetMetricsScope(ctx, project)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range got.MonitoredProjects {
		t.Logf("%s\n", v.GetName())
	}
}

func TestMetricsScopesService_CreateMonitoredProject(t *testing.T) {
	ctx := context.Background()

	project := os.Getenv("TEST_GOOGLE_CLOUD_PROJECT_MONITORED_PROJECT")

	client, err := metricsscope.NewMetricsScopesClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewMetricsScopesService(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	const monitoredProject = "hoge20180322e"
	got, err := s.CreateMonitoredProject(ctx, project, monitoredProject)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got MonitoredProject:%s\n", got.GetName())

	// すでに存在する場合はAlreadyExistsが返ってくる
	_, err = s.CreateMonitoredProject(ctx, project, monitoredProject)
	if status.Code(err) != codes.AlreadyExists {
		t.Errorf("want AlreadyExists but got %v", err)
	}

	if err := s.DeleteMonitoredProject(ctx, project, monitoredProject); err != nil {
		t.Fatal(err)
	}
}
