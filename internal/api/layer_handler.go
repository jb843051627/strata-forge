package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jb843051627/strata-forge/internal/model"
)

func (a *API) layerSubresource(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r.URL.Path, "layers")
	if !ok {
		writeError(w, model.ErrInvalidInput)
		return
	}
	if r.Method == http.MethodGet {
		item, err := a.service.GetLayer(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func pathID(path, resource string) (int64, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 4 || parts[2] != resource {
		return 0, false
	}
	id, err := strconv.ParseInt(parts[3], 10, 64)
	return id, err == nil && id > 0
}
