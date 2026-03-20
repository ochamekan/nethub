package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
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

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Unable to create connection pool", zap.Error(err))
	}
	defer dbpool.Close()

	// repo := repository.New(dbpool)
	// handler := handlers.New(repo)
	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":8000", Handler: mux}

	mux.HandleFunc("POST /devices", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /devices", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /devices/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /devices/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("DELETE /devices/{id}", func(w http.ResponseWriter, r *http.Request) {})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	var wg sync.WaitGroup
	wg.Go(func() {
		s := <-sigChan
		logger.Info("Received signal, attempting graceful shutdown", zap.Stringer("signal", s))
		cancel()
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Error("Shutdown error", zap.Error(err))
		}
		logger.Info("Gracefully stopped network devices service")
	})

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("ListenAndServe error", zap.Error(err))
	}

	wg.Wait()
}
