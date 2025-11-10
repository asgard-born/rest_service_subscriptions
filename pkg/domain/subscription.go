package domain

import (
	"context"
	"database/sql"
	"time"
)

// Subscription представляет доменную модель подписки
type Subscription struct {
	ID          string
	ServiceName string
	Price       int64
	UserID      string
	StartDate   time.Time
	EndDate     sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListFilters содержит параметры фильтрации для списка подписок
type ListFilters struct {
	UserID      string
	ServiceName string
	Limit       int
	Offset      int
}

// SummaryFilters содержит параметры фильтрации для подсчета суммы
type SummaryFilters struct {
	UserID      string
	ServiceName string
	PeriodStart time.Time
	PeriodEnd   time.Time
}

// SubscriptionRepository определяет интерфейс репозитория подписок
// Интерфейс находится в доменном слое, так как он определяет контракт для работы с доменными сущностями
type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) (*Subscription, error)
	GetByID(ctx context.Context, id string) (*Subscription, error)
	Update(ctx context.Context, id string, sub *Subscription) (*Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters ListFilters) ([]*Subscription, error)
	GetSummary(ctx context.Context, filters SummaryFilters) (int64, error)
}
