package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (context *IsoContext) video(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]

	db := context.DB

	rows, err := db.Query(`
       SELECT
		 video_key as Id,
         title as Name,
         duration as Duration,
         thumbnail_url as Thumbnail,
         url as Url
       FROM
		video
       WHERE
        video_key = ?
       LIMIT 1
    `, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var item VideoItem

	for rows.Next() {
		err := rows.Scan(&item.Id, &item.Name, &item.Duration, &item.Thumbnail, &item.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	b, err := json.Marshal(item)
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
