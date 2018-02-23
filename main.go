package main

import (
	"context"

	"github.com/nikita-vanyasin/go-web-course/common"
	"github.com/nikita-vanyasin/go-web-course/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: golint + golint in goland
	// TODO: change mysql collation
	// TODO: split context and http methods

	file := common.SetupLogging("server")
	defer file.Close()

	envSettings := common.GetEnvSettings()
	isoContext := handlers.CreateContext(envSettings)
	defer isoContext.Shutdown()

	var serverURL = ":" + envSettings.ServerPort

	log.WithFields(log.Fields{"url": serverURL}).Info("Starting the server...")

	killSignalChan := getKillSignalChan()
	srv := startServer(serverURL, isoContext)
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

func startServer(serverURL string, isoContext *handlers.IsoContext) *http.Server {
	router := handlers.Router(isoContext)
	srv := &http.Server{Addr: serverURL, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	return srv
}
