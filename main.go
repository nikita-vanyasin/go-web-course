package main

import (
	"net/http"
	"github.com/nikita-vanyasin/go-web-course/handlers"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	defer file.Close()

	const serverUrl = ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("Starting the server...")
	router := handlers.Router()
	log.Fatal(http.ListenAndServe(serverUrl, router))
}
