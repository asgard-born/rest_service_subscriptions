package usecase

import (
	"context"

	"github.com/asgard-born/rest_service_subscriptions/pkg/domain"
)

// SubscriptionRepository определяет интерфейс репозитория подписок
// Интерфейс находится в usecase, так как он используется usecase слоем
type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	GetByID(ctx context.Context, id string) (*domain.Subscription, error)
	Update(ctx context.Context, id string, sub *domain.Subscription) (*domain.Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters domain.ListFilters) ([]*domain.Subscription, error)
	GetSummary(ctx context.Context, filters domain.SummaryFilters) (int64, error)
}
