package api

import (
	"net/http"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (a *API) runSubresource(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "runs" {
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
		item, err := a.service.GetRun(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodPost:
		item, err := a.service.StartRun(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodDelete:
		item, err := a.service.CancelRun(r.Context(), id, r.URL.Query().Get("note"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
