package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/thelazydo/email-verifier/internal/config"
	"github.com/thelazydo/email-verifier/internal/database"
	"github.com/thelazydo/email-verifier/internal/middlewares"
	"github.com/thelazydo/email-verifier/internal/transport"
	"github.com/thelazydo/email-verifier/internal/worker"
	"github.com/thelazydo/email-verifier/utils"
)

func main() {
	if err := run(); err != nil {
		slog.Error("startup error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Initialize global logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	utils.LoadDisposableDomains()

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		return err
	}
	defer database.DB.Close()

	shutdownCtx, shutdownStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer shutdownStop()

	var wg sync.WaitGroup
	worker.Start(shutdownCtx, cfg.WorkersCount, &wg)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: routes(cfg),
	}

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()

	<-shutdownCtx.Done()
	slog.Info("shutdown signal recieved")

	shutdownTimeoutCtx, shutdownTimeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownTimeoutCancel()

	if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	slog.Info("waiting for workers to finish")
	wg.Wait()
	slog.Info("all workers finished")
	return nil
}

func routes(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Apply middleware chain
	statusHandler := middlewares.LimitMiddleware(transport.HandleGetStatus, cfg.RateLimit)
	VerifyHandler := middlewares.LimitMiddleware(transport.HandlePostVerify, cfg.RateLimit)

	mux.Handle("POST /verify", middlewares.LoggingMiddleware(VerifyHandler))
	mux.Handle("GET /status/{id}", middlewares.LoggingMiddleware(statusHandler))
	mux.HandleFunc("GET /", transport.HandleBaseRoute)

	return mux
}
