package handle

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("/", h.CreateSubscription)
		subscriptions.GET("/", h.ListSubscriptions)
		subscriptions.GET("/:id", h.GetSubscription)
		subscriptions.PUT("/:id", h.UpdateSubscription)
		subscriptions.DELETE("/:id", h.DeleteSubscription)
	}

	router.GET("/subscriptions/summary", h.GetSubscriptionsSummary)

	return router
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	c.JSON(200, gin.H{"message": "create subscription"})
}

func (h *Handler) ListSubscriptions(c *gin.Context) {
	c.JSON(200, gin.H{"message": "list subscriptions"})
}

func (h *Handler) GetSubscription(c *gin.Context) {
	c.JSON(200, gin.H{"message": "get subscription"})
}

func (h *Handler) UpdateSubscription(c *gin.Context) {
	c.JSON(200, gin.H{"message": "update subscription"})
}

func (h *Handler) DeleteSubscription(c *gin.Context) {
	c.JSON(200, gin.H{"message": "delete subscription"})
}

func (h *Handler) GetSubscriptionsSummary(c *gin.Context) {
	c.JSON(200, gin.H{"message": "summary of subscriptions"})
}
