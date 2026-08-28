package api

import (
	"net/http"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (a *API) reportSubresource(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "reports" {
		http.NotFound(w, r)
		return
	}
	id, ok := parsePathNumber(parts[3])
	if !ok {
		writeError(w, model.ErrInvalidInput)
		return
	}
	switch r.Method {
	case http.MethodGet:
		report, err := a.service.GetReport(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, report)
	case http.MethodPost:
		if err := a.service.ArchiveReport(r.Context(), id); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "archived"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
