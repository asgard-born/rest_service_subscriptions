package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/asgard-born/rest_service_subscriptions/pkg/domain"
	"github.com/asgard-born/rest_service_subscriptions/pkg/utils"
)

// SubscriptionUseCase содержит бизнес-логику для работы с подписками
// Реализует интерфейс SubscriptionUseCase (определен в api слое)
type SubscriptionUseCase struct {
	repo domain.SubscriptionRepository
}

// NewSubscriptionUseCase создает новый экземпляр use case для подписок
func NewSubscriptionUseCase(repo domain.SubscriptionRepository) *SubscriptionUseCase {
	return &SubscriptionUseCase{repo: repo}
}

// CreateSubscription создает новую подписку
// Реализует интерфейс api.SubscriptionUseCase
func (uc *SubscriptionUseCase) CreateSubscription(ctx context.Context, req CreateSubscriptionInput) (*domain.Subscription, error) {
	// Валидация и парсинг даты начала
	startDate, err := utils.ParseToMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	// Валидация и парсинг даты окончания
	endDate := sql.NullTime{Valid: false}
	if req.EndDate != "" {
		date, err := utils.ParseToMonthYear(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		endDate = sql.NullTime{Time: date, Valid: true}
	}

	// Валидация бизнес-правил
	if req.ServiceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Создание доменной модели
	sub := &domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       int64(req.Price),
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Сохранение через репозиторий
	created, err := uc.repo.Create(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return created, nil
}

// GetSubscription получает подписку по ID
func (uc *SubscriptionUseCase) GetSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	sub, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub, nil
}

// UpdateSubscription обновляет подписку
// Реализует интерфейс api.SubscriptionUseCase
func (uc *SubscriptionUseCase) UpdateSubscription(ctx context.Context, id string, req UpdateSubscriptionInput) (*domain.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	// Валидация и парсинг даты начала
	startDate, err := utils.ParseToMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	// Валидация и парсинг даты окончания
	endDate := sql.NullTime{Valid: false}
	if req.EndDate != "" {
		date, err := utils.ParseToMonthYear(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		endDate = sql.NullTime{Time: date, Valid: true}
	}

	// Валидация бизнес-правил
	if req.ServiceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}

	// Создание доменной модели для обновления
	sub := &domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       int64(req.Price),
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Обновление через репозиторий
	updated, err := uc.repo.Update(ctx, id, sub)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return updated, nil
}

// DeleteSubscription удаляет подписку
func (uc *SubscriptionUseCase) DeleteSubscription(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// ListSubscriptions возвращает список подписок с фильтрацией
// Реализует интерфейс api.SubscriptionUseCase
func (uc *SubscriptionUseCase) ListSubscriptions(ctx context.Context, filters ListFiltersInput) ([]*domain.Subscription, error) {
	// Валидация параметров пагинации
	if filters.Limit <= 0 {
		filters.Limit = 10 // значение по умолчанию
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	// Преобразование запроса в доменные фильтры
	domainFilters := domain.ListFilters{
		UserID:      filters.UserID,
		ServiceName: filters.ServiceName,
		Limit:       filters.Limit,
		Offset:      filters.Offset,
	}

	// Получение списка через репозиторий
	subs, err := uc.repo.List(ctx, domainFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return subs, nil
}

// GetSubscriptionsSummary вычисляет общую стоимость подписок за период
// Реализует интерфейс api.SubscriptionUseCase
func (uc *SubscriptionUseCase) GetSubscriptionsSummary(ctx context.Context, filters SummaryFiltersInput) (int64, error) {
	// Валидация обязательных полей
	if filters.PeriodStart == "" || filters.PeriodEnd == "" {
		return 0, fmt.Errorf("period_start and period_end are required")
	}

	// Парсинг дат
	periodStart, err := utils.ParseToMonthYear(filters.PeriodStart)
	if err != nil {
		return 0, fmt.Errorf("invalid period_start format: %w", err)
	}

	periodEnd, err := utils.ParseToMonthYear(filters.PeriodEnd)
	if err != nil {
		return 0, fmt.Errorf("invalid period_end format: %w", err)
	}

	// Валидация бизнес-правил
	if periodStart.After(periodEnd) {
		return 0, fmt.Errorf("period_start must be before or equal to period_end")
	}

	// Преобразование запроса в доменные фильтры
	domainFilters := domain.SummaryFilters{
		UserID:      filters.UserID,
		ServiceName: filters.ServiceName,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	}

	// Вычисление суммы через репозиторий
	total, err := uc.repo.GetSummary(ctx, domainFilters)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate summary: %w", err)
	}

	return total, nil
}

// CreateSubscriptionInput представляет входные данные для создания подписки
type CreateSubscriptionInput struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     string
}

// UpdateSubscriptionInput представляет входные данные для обновления подписки
type UpdateSubscriptionInput struct {
	ServiceName string
	Price       int
	StartDate   string
	EndDate     string
}

// ListFiltersInput представляет входные данные для получения списка подписок
type ListFiltersInput struct {
	UserID      string
	ServiceName string
	Limit       int
	Offset      int
}

// SummaryFiltersInput представляет входные данные для получения суммы подписок
type SummaryFiltersInput struct {
	UserID      string
	ServiceName string
	PeriodStart string
	PeriodEnd   string
}
