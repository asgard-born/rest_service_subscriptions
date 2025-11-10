package domain

import (
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
