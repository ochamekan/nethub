package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Starting network devices service...")

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	var wg sync.WaitGroup
	wg.Go(func() {
		s := <-sigChan
		logger.Info("Received signal, attempting graceful shutdown", zap.Stringer("signal", s))
		cancel()
		logger.Info("Gracefully stopped network devices service")
	})

	mux := http.NewServeMux()

	mux.HandleFunc("POST /devices", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /devices", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /devices/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /devices/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("DELETE /devices/{id}", func(w http.ResponseWriter, r *http.Request) {})

	logger.Fatal(http.ListenAndServe(":8000", mux).Error())

}
