package api

import (
	"context"
	"github.com/asgard-born/rest_service_subscriptions/pkg/domain"
)

// SubscriptionUseCase определяет интерфейс use case для работы с подписками
// Интерфейс определен в API слое, так как он используется здесь (dependency rule)
type SubscriptionUseCase interface {
	CreateSubscription(ctx context.Context, req CreateSubscriptionInput) (*domain.Subscription, error)
	GetSubscription(ctx context.Context, id string) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, id string, req UpdateSubscriptionInput) (*domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id string) error
	ListSubscriptions(ctx context.Context, filters ListFiltersInput) ([]*domain.Subscription, error)
	GetSubscriptionsSummary(ctx context.Context, filters SummaryFiltersInput) (int64, error)
}

// CreateSubscriptionInput представляет входные данные для создания подписки (use case слой)
type CreateSubscriptionInput struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     string
}

// UpdateSubscriptionInput представляет входные данные для обновления подписки (use case слой)
type UpdateSubscriptionInput struct {
	ServiceName string
	Price       int
	StartDate   string
	EndDate     string
}

// ListFiltersInput представляет входные данные для получения списка подписок (use case слой)
type ListFiltersInput struct {
	UserID      string
	ServiceName string
	Limit       int
	Offset      int
}

// SummaryFiltersInput представляет входные данные для получения суммы подписок (use case слой)
type SummaryFiltersInput struct {
	UserID      string
	ServiceName string
	PeriodStart string
	PeriodEnd   string
}

// SubscriptionResponse represents subscription data in API response
// swagger:model SubscriptionResponse
type SubscriptionResponse struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int64  `json:"price"`
	Currency    string `json:"currency"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
