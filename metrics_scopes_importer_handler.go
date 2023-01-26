package metrics_scope_collector

import (
	"context"
	"log"
	"net/http"
	"os"

	crmbox "github.com/sinmetalcraft/gcpbox/cloudresourcemanager/v3"
)

type MetricsScopesImporterHandler struct {
	Service *Service
}

func NewMetricsScopesImporterHandler(ctx context.Context, service *Service) (*MetricsScopesImporterHandler, error) {
	return &MetricsScopesImporterHandler{
		Service: service,
	}, nil
}

func (h *MetricsScopesImporterHandler) Handler(w http.ResponseWriter, r *http.Request) {
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
