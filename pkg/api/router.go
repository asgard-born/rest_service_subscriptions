package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
		subscriptions.GET("/summary", h.GetSubscriptionsSummary)
	}

	return router
}
