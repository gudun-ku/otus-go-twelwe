package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

)

func main() {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	log.Info("Hello world")

	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("Port is not set")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	serv := http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: router,
	}

	log.info("App is starting...")
	go serv.ListenAndServe()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	log.Info("Stopping app... ")
	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err := serv.Shutdown(timeout)
	if err != nil {
		log.Error("Error when shutdown app: %v", err)
	}
	log.Info("The app is stopped.")
}
