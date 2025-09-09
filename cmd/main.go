package main

import (
	service "github.com/asgard-born/rest_service_subscriptions"
	"github.com/asgard-born/rest_service_subscriptions/pkg/handle"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs %s", err.Error())
	}

	handler := handle.Handler{}
	srv := new(service.Server)

	if err := srv.Run(viper.GetString("8080"), handler.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
