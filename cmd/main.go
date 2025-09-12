package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "github.com/asgard-born/rest_service_subscriptions"
	_ "github.com/asgard-born/rest_service_subscriptions/docs"
	"github.com/asgard-born/rest_service_subscriptions/pkg/api"
)

// @title Subscriptions API
// @version 1.0
// @description REST API для управления подписками
// @host localhost:8080
// @BasePath /
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	slog.Info("Starting application...")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		slog.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to Postgres (pgxpool)")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("HTTP server configured", "port", port)

	srv := new(service.Server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(port, api.CreateNewRouter(pool)); err != nil {
			slog.Error("Error occurred while running HTTP server", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Warn("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited properly")
}
