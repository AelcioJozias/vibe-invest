package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AelcioJozias/vibe-invest/backend/internal/account"
	"github.com/AelcioJozias/vibe-invest/backend/internal/config"
	"github.com/AelcioJozias/vibe-invest/backend/internal/dashboard"
	"github.com/AelcioJozias/vibe-invest/backend/internal/database"
	"github.com/AelcioJozias/vibe-invest/backend/internal/investment"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	// In Spring this wiring is done mostly by DI container annotations.
	// In Go we compose dependencies explicitly, which keeps startup flow obvious.
	accountRepo := account.NewPostgresRepository(pool)
	accountService := account.NewService(accountRepo)
	accountHandler := account.NewHandler(accountService)

	investmentRepo := investment.NewPostgresRepository(pool)
	investmentService := investment.NewService(investmentRepo, accountRepo)
	investmentHandler := investment.NewHandler(investmentService)

	dashboardRepo := dashboard.NewPostgresRepository(pool)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	mux := http.NewServeMux()
	registerRoutes(mux, accountHandler, investmentHandler, dashboardHandler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           requestLogMiddleware(logger, corsMiddleware(cfg.CORSAllowOrigin, mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("server starting", "port", cfg.Port)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server failed: %w", err)
		}
	case <-sigCtx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	logger.Info("server stopped")
	return nil
}

func registerRoutes(mux *http.ServeMux, accountHandler *account.Handler, investmentHandler *investment.Handler, dashboardHandler *dashboard.Handler) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /api/v1/accounts", accountHandler.List)
	mux.HandleFunc("POST /api/v1/accounts", accountHandler.Create)
	mux.HandleFunc("PUT /api/v1/accounts/{id}", accountHandler.Update)
	mux.HandleFunc("DELETE /api/v1/accounts/{id}", accountHandler.Delete)

	mux.HandleFunc("GET /api/v1/accounts/{accountId}/investiments", investmentHandler.ListByAccount)
	mux.HandleFunc("POST /api/v1/accounts/{accountId}/investiments", investmentHandler.Create)
	mux.HandleFunc("GET /api/v1/investiments/{investmentId}", investmentHandler.GetByID)
	mux.HandleFunc("PUT /api/v1/investiments/{investmentId}", investmentHandler.Update)
	mux.HandleFunc("DELETE /api/v1/investiments/{investmentId}", investmentHandler.Delete)
	mux.HandleFunc("PUT /api/v1/investiments/{investmentId}/fees", investmentHandler.IncrementFees)

	mux.HandleFunc("GET /api/v1/portfolio/summary", dashboardHandler.Summary)
}

func requestLogMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("request handled", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
	})
}

func corsMiddleware(allowOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowOrigin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
