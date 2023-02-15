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
	parentResourceID := r.FormValue("parentResourceID")
	parentResourceType := r.FormValue("parentResourceType")

	if scopingProjectID == "" {
		log.Println("required scoping project id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if parentResourceID == "" || parentResourceType == "" {
		log.Println("required parentResourceID and parentResourceType")
		w.Write([]byte("required parentResourceID and parentResourceType"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	count, err := h.Service.ImportMonitoredProjects(ctx, scopingProjectID, &crmbox.ResourceID{ID: parentResourceID, Type: parentResourceType})
	if err != nil {
		log.Printf("failed ImportMonitoredProjects.scopingProjectID=%s,parentResourceID=%s,parentResourceType=%s. err=%s\n", scopingProjectID, parentResourceID, parentResourceType, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("finish ImportMonitoredProjects.count=%d,scopingProjectID=%s,parentResourceID=%s,parentResourceType=%s.\n", count, scopingProjectID, parentResourceID, parentResourceType)
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
