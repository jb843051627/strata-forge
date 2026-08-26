package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (a *API) samples(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		status := r.URL.Query().Get("status")
		items, err := a.service.ListSamples(r.Context(), status)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
	case http.MethodPost:
		var input model.SampleInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, model.ErrInvalidInput)
			return
		}
		item, err := a.service.ReceiveSample(r.Context(), input)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *API) sampleSubresource(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		writeError(w, model.ErrInvalidInput)
		return
	}
	if len(parts) == 4 && r.Method == http.MethodGet {
		summary, summaryErr := a.service.Summary(r.Context(), id)
		if summaryErr != nil {
			writeError(w, summaryErr)
			return
		}
		writeJSON(w, http.StatusOK, summary)
		return
	}
	if len(parts) >= 5 && parts[4] == "layers" {
		a.layersForSample(w, r, id)
		return
	}
	if len(parts) >= 5 && parts[4] == "reports" {
		reports, reportErr := a.service.ListReports(r.Context(), id)
		if reportErr != nil {
			writeError(w, reportErr)
			return
		}
		writeJSON(w, http.StatusOK, reports)
		return
	}
	if len(parts) >= 5 && parts[4] == "export" {
		data, contentType, exportErr := a.service.ExportSample(r.Context(), id, r.URL.Query().Get("format"))
		if exportErr != nil {
			writeError(w, exportErr)
			return
		}
		w.Header().Set("Content-Type", contentType)
		_, _ = w.Write(data)
		return
	}
	http.NotFound(w, r)
}

func (a *API) layersForSample(w http.ResponseWriter, r *http.Request, sampleID int64) {
	if r.Method == http.MethodGet {
		items, err := a.service.ListLayers(r.Context(), sampleID)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	if r.Method == http.MethodPost {
		var input model.LayerInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, model.ErrInvalidInput)
			return
		}
		item, err := a.service.AddLayer(r.Context(), sampleID, input)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
