package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (a *API) measurementSubresource(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "measurements" {
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
		item, err := a.service.GetMeasurement(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodPost:
		var input model.ReviewInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, model.ErrInvalidInput)
			return
		}
		review, err := a.service.ReviewMeasurement(r.Context(), id, input)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, review)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func parsePathNumber(value string) (int64, bool) {
	var id int64
	_, err := fmt.Sscan(value, &id)
	return id, err == nil && id > 0
}
