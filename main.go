package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nikita-vanyasin/go-web-course/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const serverUrl = ":8000"

func main() {
	// TODO: refactoring
	// TODO: fix toml dependecies
	// TODO: implement list params
	// TODO: retrieve path to content folder from env variable
	// TODO: remove iso context

	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	defer file.Close()

	db, err := sql.Open("mysql", `root@/simple_video_server`)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.WithFields(log.Fields{"url": serverUrl}).Info("Starting the server...")

	killSignalChan := getKillSignalChan()
	srv := startServer(serverUrl, db)
	waitForKillSignal(killSignalChan)
	srv.Shutdown(context.Background())
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("Got SIGINT")
	case syscall.SIGTERM:
		log.Info("Got SIGTERM")
	}
}

func startServer(serverUrl string, db *sql.DB) *http.Server {
	router := handlers.Router(db)
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	return srv
}
