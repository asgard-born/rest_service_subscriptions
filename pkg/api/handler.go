package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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

func CreateNewRouter(db *pgxpool.Pool) *gin.Engine {
	h := Handler{
		db: db,
	}

	router := gin.New()

	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("/", h.CreateSubscription)
		subscriptions.PUT("/:id", h.UpdateSubscription)
		subscriptions.GET("/:id", h.GetSubscription)
		subscriptions.GET("/", h.ListSubscriptions)
		subscriptions.DELETE("/:id", h.DeleteSubscription)
	}

	router.GET("/subscriptions/summary", h.GetSubscriptionsSummary)

	return router
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest

	if err := c.BindJSON(&req); err != nil {
		log.Println(fmt.Sprintf("err ocurred while binding %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is epmty"})
		return
	}

	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user-id is epmty"})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date, use MM-YYYY"})
		return
	}

	var endDate *time.Time

	if req.EndDate != nil {
		date, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date, use MM-YYYY"})
			return
		}
		endDate = &date
	}

	err = h.db.QueryRow(context.Background(),
		"INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		req.ServiceName, req.Price, req.UserID, startDate, endDate).Scan(&id)

	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}
}

func (h *Handler) UpdateSubscription(c *gin.Context) {

}

func (h *Handler) GetSubscription(c *gin.Context) {
}

func (h *Handler) ListSubscriptions(c *gin.Context) {

}

func (h *Handler) DeleteSubscription(c *gin.Context) {
}

func (h *Handler) GetSubscriptionsSummary(c *gin.Context) {
}
