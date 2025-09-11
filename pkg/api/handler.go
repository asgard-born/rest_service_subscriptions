package api

import (
	"github.com/asgard-born/rest_service_subscriptions/pkg/db"
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

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	slog.Info("CreateSubscription called")

	var req CreateSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		slog.Warn("Invalid start_date", "value", req.StartDate)
		RespondError(c, http.StatusBadRequest, "invalid start_date, use MM-YYYY")
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		date, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			slog.Warn("Invalid end_date", "value", *req.EndDate)
			RespondError(c, http.StatusBadRequest, "invalid end_date, use MM-YYYY")
			return
		}
		endDate = &date
	}

	var id string
	err = h.db.QueryRow(
		c.Request.Context(),
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id`,
		req.ServiceName, req.Price, req.UserID, startDate, endDate,
	).Scan(&id)

	if err != nil {
		slog.Error("Failed to insert subscription", "error", err)
		RespondError(c, http.StatusInternalServerError, "failed to insert subscription")
		return
	}

	slog.Info("Subscription created", "id", id)
	RespondSuccess(c, http.StatusCreated, gin.H{
		"id":      id,
		"message": "subscription created successfully",
	})
}

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
		c,
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

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		slog.Warn("Invalid start_date", "value", req.StartDate)
		RespondError(c, http.StatusBadRequest, "invalid start_date, use MM-YYYY")
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		date, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			slog.Warn("Invalid end_date", "value", *req.EndDate)
			RespondError(c, http.StatusBadRequest, "invalid end_date, use MM-YYYY")
			return
		}
		endDate = &date
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
	RespondSuccess(c, http.StatusOK, sub)
}

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

func (h *Handler) ListSubscriptions(c *gin.Context) {

}

func (h *Handler) GetSubscriptionsSummary(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	periodStart := c.Query("period_start")
	periodEnd := c.Query("period_end")

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE start_date >= $1 AND (end_date IS NULL OR end_date <= $2)
	`
	args := []interface{}{periodStart, periodEnd}

	if userID != "" {
		query += " AND user_id = $3"
		args = append(args, userID)
	}
	if serviceName != "" {
		query += " AND service_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, serviceName)
	}

	var total float64

	err := h.db.QueryRow(c.Request.Context(), query, args...).Scan(&total)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, "failed to calculate summary")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"timestamp": time.Now().UTC(),
	})
}
