package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"

	service "github.com/asgard-born/rest_service_subscriptions"
	"github.com/asgard-born/rest_service_subscriptions/pkg/api"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	dbpool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping db: %s", err.Error())
	}

	log.Println("✅ Connected to Postgres (pgxpool)")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := new(service.Server)

	if err := srv.Run(port, api.CreateNewRouter(dbpool)); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
