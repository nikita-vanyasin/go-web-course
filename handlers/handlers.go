package handlers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/nikita-vanyasin/go-web-course/video"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type IsoContext struct {
	DB              *sql.DB
	VideoRepository video.RepositoryInterface
}

func Router(db *sql.DB) http.Handler {

	repo := video.CreateRepository(db)
	context := IsoContext{db, repo}

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/list", context.list).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", context.video).Methods(http.MethodGet)
	s.HandleFunc("/video", context.postVideo).Methods(http.MethodPost)
	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL.Path,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("Got a new request")
		h.ServeHTTP(w, r)
	})
}
