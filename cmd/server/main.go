package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	metricsscope "cloud.google.com/go/monitoring/metricsscope/apiv1"
	crmbox "github.com/sinmetalcraft/gcpbox/cloudresourcemanager/v3"
	metricsscopebox "github.com/sinmetalcraft/gcpbox/monitoring/metricsscope/v0"
	msc "github.com/sinmetalcraft/metrics-scope-collector"
	"google.golang.org/api/cloudresourcemanager/v3"
)

func main() {
	ctx := context.Background()

	log.Print("starting server...")
	client, err := metricsscope.NewMetricsScopesClient(ctx)
	if err != nil {
		panic(err)
	}

	metricsScopesService, err := metricsscopebox.NewService(ctx, client)
	if err != nil {
		panic(err)
	}

	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		panic(err)
	}

	resourceManagerService, err := crmbox.NewResourceManagerService(ctx, crmService)
	if err != nil {
		panic(err)
	}

	s, err := msc.NewService(ctx, metricsScopesService, resourceManagerService)
	if err != nil {
		panic(err)
	}

	metricsScopesImporterHandler, err := msc.NewMetricsScopesImporterHandler(ctx, s)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/metrics-scope-gatherer/create", metricsScopesImporterHandler.CreateHandler)
	http.HandleFunc("/metrics-scope-gatherer/cleanup", metricsScopesImporterHandler.CleanUpHandler)

	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\n")
}
