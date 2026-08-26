package api

import (
	"net/http"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (a *API) alertSubresource(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "alerts" {
		http.NotFound(w, r)
		return
	}
	id, ok := parsePathNumber(parts[3])
	if !ok {
		writeError(w, model.ErrInvalidInput)
		return
	}
	if r.Method == http.MethodPost {
		if err := a.service.AcknowledgeAlert(r.Context(), id); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "acknowledged"})
		return
	}
	if r.Method == http.MethodGet {
		item, err := a.service.Store().GetAlert(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
