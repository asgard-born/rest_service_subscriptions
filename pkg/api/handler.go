package api

import (
	"context"
	"database/sql"
	"github.com/asgard-born/rest_service_subscriptions/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	db *pgxpool.Pool
}

type CreateSubscriptionRequest struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
}

type UpdateSubscriptionRequest struct {
	ServiceName string
	Price       int
	StartDate   string
	EndDate     *string
}

type ErrorResponse struct {
	Error     string `json:"error"`
	Details   string `json:"details,omitempty"`
	Code      int    `json:"code"`
	Timestamp string `json:"timestamp"`
}

func CreateNewRouter(db *pgxpool.Pool) *gin.Engine {
	h := Handler{
		db: db,
	}

	router := gin.New()

	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("/", h.CreateSubscription)
		subscriptions.GET("/:id", h.GetSubscription)
		subscriptions.PUT("/:id", h.UpdateSubscription)
		subscriptions.DELETE("/:id", h.DeleteSubscription)
		subscriptions.GET("/", h.ListSubscriptions)
	}

	router.GET("/subscriptions/summary", h.GetSubscriptionsSummary)

	return router
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	log.Println("CreateSubscription called")

	var req CreateSubscriptionRequest

	if err := c.BindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON: %v", err)

		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "invalid request body",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	log.Printf("Request: %+v", req)

	if _, err := uuid.Parse(req.UserID); err != nil {
		log.Printf("Invalid user_id: %s", req.UserID)

		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "invalid user_id, must be UUID",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		log.Printf("Invalid start_date: %s", req.StartDate)

		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "invalid start_date, use MM-YYYY",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		date, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			log.Printf("Invalid end_date: %s", *req.EndDate)

			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:     "invalid end_date, use MM-YYYY",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
		endDate = &date
	}

	var id string
	err = h.db.QueryRow(
		context.Background(),
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		req.ServiceName, req.Price, req.UserID, startDate, endDate,
	).Scan(&id)

	if err != nil {
		log.Printf("Failed to insert subscription: %v", err)

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:     "failed to insert subscription",
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	log.Printf("Subscription created with id=%s", id)

	c.JSON(http.StatusCreated, gin.H{
		"id":        id,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) GetSubscription(c *gin.Context) {
	log.Println("GetSubscription called")
	id := c.Param("id")

	if id == "" {
		log.Println("missing id param")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "id is required",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	log.Printf("Getting subscription id=%s", id)

	var sub repository.Subscription

	err := h.db.QueryRow(
		context.Background(),
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

	if err == sql.ErrNoRows {
		log.Printf("Subscription not found: id=%s", id)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:     "subscription not found",
			Code:      http.StatusNotFound,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}
	if err != nil {
		log.Printf("Failed to get subscription id=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:     "failed to get subscription",
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	log.Printf("Subscription retrieved: %+v", sub)
	c.JSON(http.StatusOK, sub)
}

func (h *Handler) UpdateSubscription(c *gin.Context) {
	log.Println("UpdateSubscription called")

	id := c.Param("id")

	if id == "" {
		log.Println("missing id param")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "id is required",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	log.Printf("Updating subscription id=%s", id)

	var req UpdateSubscriptionRequest

	if err := c.BindJSON(&req); err != nil {
		log.Printf("failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "invalid request body",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	log.Printf("Update request: %+v", req)

	startDate, err := time.Parse("01-2006", req.StartDate)

	if err != nil {
		log.Printf("Invalid start_date: %s", req.StartDate)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:     "invalid start_date, use MM-YYYY",
			Code:      http.StatusBadRequest,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	var endDate *time.Time

	if req.EndDate != nil {
		date, err := time.Parse("01-2006", *req.EndDate)

		if err != nil {
			log.Printf("Invalid end_date: %s", *req.EndDate)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:     "invalid end_date, use MM-YYYY",
				Code:      http.StatusBadRequest,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
		endDate = &date
	}

	var sub repository.Subscription

	err = h.db.QueryRow(
		context.Background(),
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

	if err == sql.ErrNoRows {
		log.Printf("Subscription not found: id=%s", id)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:     "subscription not found",
			Code:      http.StatusNotFound,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}
	if err != nil {
		log.Printf("Failed to update subscription id=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:     "failed to update subscription",
			Code:      http.StatusInternalServerError,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	log.Printf("Subscription updated: %+v", sub)
	c.JSON(http.StatusOK, sub)
}

func (h *Handler) DeleteSubscription(c *gin.Context) {
}

func (h *Handler) ListSubscriptions(c *gin.Context) {

}

func (h *Handler) GetSubscriptionsSummary(c *gin.Context) {
}
