package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

const defaultLimit = 10000

func getParam(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func getIntParam(r *http.Request, name string) uint64 {
	result := uint64(0)
	if param := getParam(r, name); param != "" {
		result, _ = strconv.ParseUint(param, 10, 64)
	}
	return result
}

func (context *IsoContext) list(w http.ResponseWriter, r *http.Request) {
	searchStringParam := getParam(r, "searchString")
	skip := getIntParam(r, "skip")
	limit := getIntParam(r, "limit")
	if limit < 1 {
		limit = defaultLimit
	}
	list, err := context.VideoRepository.List(searchStringParam, skip, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("error writing response")
	}
}
