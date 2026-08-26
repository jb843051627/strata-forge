package api

import (
	"net/http"

	"github.com/jb843051627/strata-forge/internal/report"
	"github.com/jb843051627/strata-forge/internal/service"
)

type API struct {
	service  *service.LabService
	renderer *report.Renderer
}

func NewHandler(svc *service.LabService, renderer *report.Renderer) http.Handler {
	api := &API{service: svc, renderer: renderer}
	mux := http.NewServeMux()
	mux.HandleFunc("/", api.page)
	mux.HandleFunc("/healthz", api.health)
	mux.HandleFunc("/api/v1/search", api.search)
	mux.HandleFunc("/api/v1/metrics", api.metrics)
	mux.HandleFunc("/api/v1/events/", api.eventSubresource)
	mux.HandleFunc("/api/v1/samples", api.samples)
	mux.HandleFunc("/api/v1/samples/", api.sampleSubresource)
	mux.HandleFunc("/api/v1/layers/", api.layerSubresource)
	mux.HandleFunc("/api/v1/measurements/", api.measurementSubresource)
	mux.HandleFunc("/api/v1/reviews/", api.reviewSubresource)
	mux.HandleFunc("/api/v1/runs/", api.runSubresource)
	mux.HandleFunc("/api/v1/reports/", api.reportSubresource)
	mux.HandleFunc("/api/v1/alerts/", api.alertSubresource)
	return requestLog(mux)
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}
