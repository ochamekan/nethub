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
	"github.com/ochamekan/nethub/backend/internal/handlers"
	"github.com/ochamekan/nethub/backend/internal/repository"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("starting network devices service...")

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("error loading .env file", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("unable to create connection pool", zap.Error(err))
	}
	defer dbpool.Close()

	repo := repository.New(dbpool)
	handler := handlers.New(repo, logger)
	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":8000", Handler: mux}

	mux.HandleFunc("POST /devices", handler.CreateDevice)
	mux.HandleFunc("GET /devices", handler.GetDevices)
	mux.HandleFunc("GET /devices/{id}", handler.GetDevice)
	mux.HandleFunc("POST /devices/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("DELETE /devices/{id}", func(w http.ResponseWriter, r *http.Request) {})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	var wg sync.WaitGroup
	wg.Go(func() {
		s := <-sigChan
		logger.Info("received signal, attempting graceful shutdown", zap.Stringer("signal", s))
		cancel()
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Error("shutdown error", zap.Error(err))
		}
		logger.Info("gracefully stopped network devices service")
	})

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("ListenAndServe error", zap.Error(err))
	}

	wg.Wait()
}
