package api

import (
	"database/sql"
	"fmt"
	"github.com/asgard-born/rest_service_subscriptions/pkg/db"
	"github.com/asgard-born/rest_service_subscriptions/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	db *pgxpool.Pool
}

// CreateSubscriptionRequest represents data for creating a subscription
// swagger:model CreateSubscriptionRequest
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

// UpdateSubscriptionRequest represents data for updating a subscription
// swagger:model UpdateSubscriptionRequest
type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

// CreateSubscription godoc
// @Summary Создать подписку
// @Description Создаёт новую подписку для пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	slog.Info("CreateSubscription called")

	var req CreateSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	startDate, err := utils.ParseToMonthYear(req.StartDate)
	if err != nil {
		slog.Warn("Invalid start_date", "value", req.StartDate, "err", err)
		RespondError(c, http.StatusBadRequest, "invalid start_date format")
		return
	}

	endDate := sql.NullTime{Valid: false}
	if req.EndDate != "" {
		date, err := utils.ParseToMonthYear(req.EndDate)
		if err != nil {
			slog.Warn("Invalid end_date", "value", req.EndDate, "err", err)
			RespondError(c, http.StatusBadRequest, "invalid end_date")
			return
		}
		endDate = sql.NullTime{Time: date, Valid: true}
	}

	var sub db.Subscription
	err = h.db.QueryRow(
		c.Request.Context(),
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`,
		req.ServiceName, req.Price, req.UserID, startDate, endDate,
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

	if err != nil {
		slog.Error("Failed to insert subscription", "error", err)
		RespondError(c, http.StatusInternalServerError, "failed to insert subscription")
		return
	}

	slog.Info("Subscription created", "id", sub.ID)
	RespondSuccess(c, http.StatusCreated, ToSubscriptionResponse(sub))
}

// GetSubscription godoc
// @Summary Получить подписку
// @Description Возвращает подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	slog.Info("GetSubscription called")

	id := c.Param("id")
	if id == "" {
		slog.Warn("Missing id param")
		RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	var sub db.Subscription
	err := h.db.QueryRow(
		c.Request.Context(),
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
		&sub.UpdatedAt)

	if err == pgx.ErrNoRows {
		slog.Warn("Subscription not found", "id", id)
		RespondError(c, http.StatusNotFound, "subscription not found")
		return
	}
	if err != nil {
		slog.Error("Failed to get subscription", "id", id, "error", err)
		RespondError(c, http.StatusInternalServerError, "failed to get subscription")
		return
	}

	slog.Info("Subscription retrieved", "id", id)
	RespondSuccess(c, http.StatusOK, ToSubscriptionResponse(sub))
}

// UpdateSubscription godoc
// @Summary Обновить подписку
// @Description Обновляет существующую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	slog.Info("UpdateSubscription called")

	id := c.Param("id")
	if id == "" {
		slog.Warn("Missing id param")
		RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	startDate, err := utils.ParseToMonthYear(req.StartDate)
	if err != nil {
		slog.Warn("Invalid start_date", "value", req.StartDate, "err", err)
		RespondError(c, http.StatusBadRequest, "invalid start_date format")
		return
	}

	endDate := sql.NullTime{Valid: false}
	if req.EndDate != "" {
		date, err := utils.ParseToMonthYear(req.EndDate)
		if err != nil {
			slog.Warn("Invalid end_date", "value", req.EndDate, "err", err)
			RespondError(c, http.StatusBadRequest, "invalid end_date")
			return
		}
		endDate = sql.NullTime{Time: date, Valid: true}
	}

	var sub db.Subscription

	err = h.db.QueryRow(
		c.Request.Context(),
		`UPDATE subscriptions
         SET service_name = $1,
             price = $2,
             start_date = $3,
             end_date = $4,
             updated_at = now()
         WHERE id = $5
         RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`,
		req.ServiceName, req.Price, startDate, endDate, id,
	).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt)

	if err == pgx.ErrNoRows {
		slog.Warn("Subscription not found for update", "id", id)
		RespondError(c, http.StatusNotFound, "subscription not found")
		return
	}
	if err != nil {
		slog.Error("Failed to update subscription", "id", id, "error", err)
		RespondError(c, http.StatusInternalServerError, "failed to update subscription")
		return
	}

	slog.Info("Subscription updated", "id", id)
	RespondSuccess(c, http.StatusOK, ToSubscriptionResponse(sub))
}

// DeleteSubscription godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	slog.Info("DeleteSubscription called")

	id := c.Param("id")
	if id == "" {
		slog.Warn("Missing id param")
		RespondError(c, http.StatusBadRequest, "id is required")
		return
	}

	cmdTag, err := h.db.Exec(
		c.Request.Context(),
		`DELETE FROM subscriptions WHERE id = $1`,
		id,
	)

	if err != nil {
		slog.Error("Failed to delete subscription", "id", id, "error", err)
		RespondError(c, http.StatusInternalServerError, "failed to delete subscription")
		return
	}

	if cmdTag.RowsAffected() == 0 {
		slog.Warn("Subscription not found for delete", "id", id)
		RespondError(c, http.StatusNotFound, "subscription not found")
		return
	}

	slog.Info("Subscription deleted", "id", id)
	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "subscription deleted successfully",
	})
}

// ListSubscriptions godoc
// @Summary Список подписок
// @Description Возвращает список подписок с пагинацией и фильтрацией
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "Фильтр по user_id"
// @Param service_name query string false "Фильтр по service_name"
// @Param limit query int false "Лимит (по умолчанию 10)" default(10)
// @Param offset query int false "Смещение (по умолчанию 0)" default(0)
// @Success 200 {array} SubscriptionResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	slog.Info("ListSubscriptions called",
		"user_id", userID,
		"service_name", serviceName,
		"limit", limitStr,
		"offset", offsetStr,
	)

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		slog.Warn("invalid limit", "value", limitStr, "err", err)
		RespondError(c, http.StatusBadRequest, "invalid limit")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		slog.Warn("invalid offset", "value", offsetStr, "err", err)
		RespondError(c, http.StatusBadRequest, "invalid offset")
		return
	}

	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
			  FROM subscriptions
			  WHERE 1=1`

	var args []interface{}
	argIndex := 1

	if userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}
	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, serviceName)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(c.Request.Context(), query, args...)
	if err != nil {
		slog.Error("failed to query subscriptions", "err", err)
		RespondError(c, http.StatusInternalServerError, "failed to list subscriptions")
		return
	}
	defer rows.Close()

	var subs []SubscriptionResponse
	for rows.Next() {
		var s db.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
			slog.Error("failed to scan subscription", "err", err)
			RespondError(c, http.StatusInternalServerError, "failed to scan subscription")
			return
		}
		subs = append(subs, ToSubscriptionResponse(s))
	}

	slog.Info("subscriptions listed", "count", len(subs))
	RespondSuccess(c, http.StatusOK, subs)
}

// GetSubscriptionsSummary godoc
// @Summary Сумма подписок
// @Description Возвращает общую стоимость подписок за период
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "Фильтр по user_id"
// @Param service_name query string false "Фильтр по service_name"
// @Param period_start query string true "Начало периода (MM-YYYY)"
// @Param period_end query string true "Конец периода (MM-YYYY)"
// @Success 200 {object} object{total=int64,from=string,to=string,user_id=string,service=string,timestamp=string}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /subscriptions/summary [get]
func (h *Handler) GetSubscriptionsSummary(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	periodStartQuery := c.Query("period_start")
	periodEndQuery := c.Query("period_end")

	slog.Info("GetSubscriptionsSummary called",
		"user_id", userID,
		"service_name", serviceName,
		"period_start", periodStartQuery,
		"period_end", periodEndQuery,
	)

	if periodStartQuery == "" || periodEndQuery == "" {
		RespondError(c, http.StatusBadRequest, "period_start and period_end are required")
		return
	}

	periodStart, err := utils.ParseToMonthYear(periodStartQuery)
	if err != nil {
		slog.Warn("invalid period_start", "value", periodStartQuery, "err", err)
		RespondError(c, http.StatusBadRequest, "invalid period_start")
		return
	}

	periodEnd, err := utils.ParseToMonthYear(periodEndQuery)
	if err != nil {
		slog.Warn("invalid period_end", "value", periodEndQuery, "err", err)
		RespondError(c, http.StatusBadRequest, "invalid period_end")
		return
	}

	query :=
		`SELECT COALESCE(SUM(
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

	args := []interface{}{periodStart, periodEnd}

	if userID != "" {
		query += " AND user_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, userID)
	}
	if serviceName != "" {
		query += " AND service_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, serviceName)
	}

	var total int64
	err = h.db.QueryRow(c.Request.Context(), query, args...).Scan(&total)
	if err != nil {
		slog.Error("failed to calculate summary", "err", err)
		RespondError(c, http.StatusInternalServerError, "failed to calculate summary")
		return
	}

	slog.Info("summary calculated",
		"total", total,
		"from", periodStart,
		"to", periodEnd,
	)

	RespondSuccess(c, http.StatusOK, gin.H{
		"total":     total,
		"from":      periodStart.Format("01-2006"),
		"to":        periodEnd.Format("01-2006"),
		"user_id":   userID,
		"service":   serviceName,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
