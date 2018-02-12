package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

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
	/*
		searchStringParam := getParam(r, "searchString")
		skip := getIntParam(r, "skip")
		limit := getIntParam(r, "limit")
	*/
	db := context.DB

	rows, err := db.Query(`
       SELECT
		 video_key as Id,
         title as Name,
         duration as Duration,
         thumbnail_url as Thumbnail,
         url as Url
       FROM video
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []VideoItem
	for rows.Next() {
		var item VideoItem
		err := rows.Scan(&item.Id, &item.Name, &item.Duration, &item.Thumbnail, &item.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, item)
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
