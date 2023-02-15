package metrics_scope_collector

import (
	"context"
	"log"
	"net/http"
	"os"

	crmbox "github.com/sinmetalcraft/gcpbox/cloudresourcemanager/v3"
)

type MetricsScopesGathererHandler struct {
	Service *Service
}

func NewMetricsScopesImporterHandler(ctx context.Context, service *Service) (*MetricsScopesGathererHandler, error) {
	return &MetricsScopesGathererHandler{
		Service: service,
	}, nil
}

func (h *MetricsScopesGathererHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	scopingProjectID := os.Getenv("SCOPING_PROJECT_ID")
	importResourceID := r.FormValue("importResourceID")
	importResourceType := r.FormValue("importResourceType")

	count, err := h.Service.ImportMonitoredProjects(ctx, scopingProjectID, &crmbox.ResourceID{ID: importResourceID, Type: importResourceType})
	if err != nil {
		log.Printf("failed ImportMonitoredProjects.scopingProjectID=%s,importResourceID=%s,importResourceType=%s. err=%s\n", scopingProjectID, importResourceID, importResourceType, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("finish ImportMonitoredProjects.count=%d,scopingProjectID=%s,importResourceID=%s,importResourceType=%s.\n", count, scopingProjectID, importResourceID, importResourceType)
}

func (h *MetricsScopesGathererHandler) CleanUpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	scopingProjectID := os.Getenv("SCOPING_PROJECT_ID")
	count, err := h.Service.CleanUp(ctx, scopingProjectID)
	if err != nil {
		log.Printf("failed CleanUp.count=%d,scopingProjectID=%s err=%s\n", count, scopingProjectID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("finish CleanUpMonitoredProjects.count=%d,scopingProjectID=%s\n", count, scopingProjectID)
}
