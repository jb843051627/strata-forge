package api

import (
	"net/http"
	"strconv"
)

func writeLocation(w http.ResponseWriter, path string, id int64) {
	w.Header().Set("Location", path+"/"+strconv.FormatInt(id, 10))
}

func accepted(w http.ResponseWriter, value any) {
	writeJSON(w, http.StatusAccepted, value)
}
