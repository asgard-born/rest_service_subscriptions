package main

import (
	"database/sql"
	"log"
	"os"

	service "github.com/asgard-born/rest_service_subscriptions"
	"github.com/asgard-born/rest_service_subscriptions/pkg/handle"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// читаем ENV
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// подключение к БД
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %s", err.Error())
	}

	log.Println("✅ Connected to Postgres")

	// запуск сервера
	handler := handle.Handler{}
	srv := new(service.Server)

	if err := srv.Run(port, handler.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
