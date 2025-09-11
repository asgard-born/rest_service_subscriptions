package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "github.com/asgard-born/rest_service_subscriptions"
	"github.com/asgard-born/rest_service_subscriptions/pkg/api"
)

func main() {
	log.Println("Starting application...")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to Postgres (pgxpool)")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("HTTP server will listen on port %s", port)

	srv := new(service.Server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(port, api.CreateNewRouter(pool)); err != nil {
			log.Fatalf("Error occurred while running HTTP server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
