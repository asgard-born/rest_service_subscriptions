package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	service "github.com/asgard-born/rest_service_subscriptions"
	_ "github.com/asgard-born/rest_service_subscriptions/docs"
	"github.com/asgard-born/rest_service_subscriptions/pkg/api"
	"github.com/asgard-born/rest_service_subscriptions/pkg/domain"
	"github.com/asgard-born/rest_service_subscriptions/pkg/infrastructure/postgres"
	"github.com/asgard-born/rest_service_subscriptions/pkg/usecase"
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

	// Инициализация базы данных
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

	// Инициализация слоев архитектуры
	// Infrastructure layer (инфраструктурный слой, реализует доменные интерфейсы)
	subscriptionRepo := postgres.NewSubscriptionRepository(pool)

	// UseCase layer (бизнес-логика)
	subscriptionUseCaseImpl := usecase.NewSubscriptionUseCase(subscriptionRepo)

	// Адаптер для преобразования типов между API и UseCase слоями
	subscriptionUseCase := &useCaseAdapter{uc: subscriptionUseCaseImpl}

	// API layer (хэндлеры и роутер)
	router := api.CreateNewRouter(subscriptionUseCase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("HTTP server configured", "port", port)

	srv := new(service.Server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(port, router); err != nil {
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

// useCaseAdapter адаптирует usecase.SubscriptionUseCase к интерфейсу api.SubscriptionUseCase
// Преобразует типы между API и UseCase слоями
type useCaseAdapter struct {
	uc *usecase.SubscriptionUseCase
}

func (a *useCaseAdapter) CreateSubscription(ctx context.Context, req api.CreateSubscriptionInput) (*domain.Subscription, error) {
	ucReq := usecase.CreateSubscriptionInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
	return a.uc.CreateSubscription(ctx, ucReq)
}

func (a *useCaseAdapter) GetSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	return a.uc.GetSubscription(ctx, id)
}

func (a *useCaseAdapter) UpdateSubscription(ctx context.Context, id string, req api.UpdateSubscriptionInput) (*domain.Subscription, error) {
	ucReq := usecase.UpdateSubscriptionInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
	return a.uc.UpdateSubscription(ctx, id, ucReq)
}

func (a *useCaseAdapter) DeleteSubscription(ctx context.Context, id string) error {
	return a.uc.DeleteSubscription(ctx, id)
}

func (a *useCaseAdapter) ListSubscriptions(ctx context.Context, filters api.ListFiltersInput) ([]*domain.Subscription, error) {
	ucFilters := usecase.ListFiltersInput{
		UserID:      filters.UserID,
		ServiceName: filters.ServiceName,
		Limit:       filters.Limit,
		Offset:      filters.Offset,
	}
	return a.uc.ListSubscriptions(ctx, ucFilters)
}

func (a *useCaseAdapter) GetSubscriptionsSummary(ctx context.Context, filters api.SummaryFiltersInput) (int64, error) {
	ucFilters := usecase.SummaryFiltersInput{
		UserID:      filters.UserID,
		ServiceName: filters.ServiceName,
		PeriodStart: filters.PeriodStart,
		PeriodEnd:   filters.PeriodEnd,
	}
	return a.uc.GetSubscriptionsSummary(ctx, ucFilters)
}
