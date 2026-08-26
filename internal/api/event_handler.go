package api

import (
	"net/http"
	"strings"
)

func (a *API) eventSubresource(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "events" {
		http.NotFound(w, r)
		return
	}
	id, ok := parsePathNumber(parts[3])
	if !ok {
		http.Error(w, "invalid sample id", http.StatusBadRequest)
		return
	}
	events, err := a.service.Store().ListEvents(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, events)
}
