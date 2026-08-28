package api

import (
	"net/http"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (a *API) reviewSubresource(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "reviews" {
		http.NotFound(w, r)
		return
	}
	id, ok := parsePathNumber(parts[3])
	if !ok {
		writeError(w, model.ErrInvalidInput)
		return
	}
	item, err := a.service.LatestReview(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
