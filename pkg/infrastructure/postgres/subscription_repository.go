package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/asgard-born/rest_service_subscriptions/pkg/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Проверка, что SubscriptionRepository реализует интерфейс domain.SubscriptionRepository
var _ domain.SubscriptionRepository = (*SubscriptionRepository)(nil)

// SubscriptionRepository реализует интерфейс репозитория для PostgreSQL
type SubscriptionRepository struct {
	db *pgxpool.Pool
}

// NewSubscriptionRepository создает новый экземпляр репозитория подписок
func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create создает новую подписку
func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	var created domain.Subscription
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(
		&created.ID,
		&created.ServiceName,
		&created.Price,
		&created.UserID,
		&created.StartDate,
		&created.EndDate,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return &created, nil
}

// GetByID получает подписку по ID
func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.db.QueryRow(
		ctx,
		`SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at 
         FROM subscriptions 
         WHERE id = $1`,
		id,
	).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &sub, nil
}

// Update обновляет подписку
func (r *SubscriptionRepository) Update(ctx context.Context, id string, sub *domain.Subscription) (*domain.Subscription, error) {
	var updated domain.Subscription
	err := r.db.QueryRow(
		ctx,
		`UPDATE subscriptions
         SET service_name = $1,
             price = $2,
             start_date = $3,
             end_date = $4,
             updated_at = now()
         WHERE id = $5
         RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`,
		sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, id,
	).Scan(
		&updated.ID,
		&updated.ServiceName,
		&updated.Price,
		&updated.UserID,
		&updated.StartDate,
		&updated.EndDate,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return &updated, nil
}

// Delete удаляет подписку
func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.db.Exec(
		ctx,
		`DELETE FROM subscriptions WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// List возвращает список подписок с фильтрацией
func (r *SubscriptionRepository) List(ctx context.Context, filters domain.ListFilters) ([]*domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
			  FROM subscriptions
			  WHERE 1=1`

	var args []interface{}
	argIndex := 1

	if filters.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, filters.UserID)
		argIndex++
	}
	if filters.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, filters.ServiceName)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []*domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return subs, nil
}

// GetSummary вычисляет общую стоимость подписок за период
func (r *SubscriptionRepository) GetSummary(ctx context.Context, filters domain.SummaryFilters) (int64, error) {
	query := `SELECT COALESCE(SUM(
		  CASE WHEN start_date <= $2 AND COALESCE(end_date, $2) >= $1
			THEN price * (
			  (EXTRACT(YEAR FROM LEAST(COALESCE(end_date, $2), $2)) - EXTRACT(YEAR FROM GREATEST(start_date, $1))) * 12 +
			  EXTRACT(MONTH FROM LEAST(COALESCE(end_date, $2), $2)) - EXTRACT(MONTH FROM GREATEST(start_date, $1)) + 1
			)
			ELSE 0
		  END
		), 0) AS total
		FROM subscriptions
		WHERE 1=1`

	args := []interface{}{filters.PeriodStart, filters.PeriodEnd}

	if filters.UserID != "" {
		query += " AND user_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, filters.UserID)
	}
	if filters.ServiceName != "" {
		query += " AND service_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, filters.ServiceName)
	}

	var total int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate summary: %w", err)
	}

	return total, nil
}
