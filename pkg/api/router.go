package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// CreateNewRouter создает новый роутер с инициализированными хэндлерами
func CreateNewRouter(subscriptionUseCase SubscriptionUseCase) *gin.Engine {
	h := NewHandler(subscriptionUseCase)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("/", h.CreateSubscription)
		subscriptions.GET("/:id", h.GetSubscription)
		subscriptions.PUT("/:id", h.UpdateSubscription)
		subscriptions.DELETE("/:id", h.DeleteSubscription)
		subscriptions.GET("/", h.ListSubscriptions)
		subscriptions.GET("/summary", h.GetSubscriptionsSummary)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
