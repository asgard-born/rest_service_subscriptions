package main

import (
	"database/sql"
	"log"

	service "github.com/asgard-born/rest_service_subscriptions"
	"github.com/asgard-born/rest_service_subscriptions/pkg/handle"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	dsn := viper.GetString("db.dsn")
	if dsn == "" {
		dsn = getenv("DATABASE_URL", "")
	}

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatalf("failed to open db: %s", err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %s", err.Error())
	}

	log.Println("✅ Connected to Postgres")

	handler := handle.Handler{}
	srv := new(service.Server)

	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}

	defer db.Close()
}

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}

func getenv(key, fallback string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}

	return fallback
}
