package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	subscriptionUseCase SubscriptionUseCase
}

// NewHandler создает новый экземпляр хэндлера
func NewHandler(subscriptionUseCase SubscriptionUseCase) *Handler {
	return &Handler{
		subscriptionUseCase: subscriptionUseCase,
	}
}

// CreateSubscriptionRequest represents data for creating a subscription
// swagger:model CreateSubscriptionRequest
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,min=0"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date,omitempty"`
}

// UpdateSubscriptionRequest represents data for updating a subscription
// swagger:model UpdateSubscriptionRequest
type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,min=0"`
	StartDate   string `json:"start_date" binding:"required"`
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
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Failed to bind JSON", "error", err)
		RespondError(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// Преобразование HTTP запроса в use case запрос
	useCaseReq := CreateSubscriptionInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	// Вызов use case
	sub, err := h.subscriptionUseCase.CreateSubscription(c.Request.Context(), useCaseReq)
	if err != nil {
		slog.Error("Failed to create subscription", "error", err)
		handleError(c, err)
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

	// Вызов use case
	sub, err := h.subscriptionUseCase.GetSubscription(c.Request.Context(), id)
	if err != nil {
		slog.Error("Failed to get subscription", "id", id, "error", err)
		handleError(c, err)
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
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Failed to bind JSON", "error", err)
		RespondError(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// Преобразование HTTP запроса в use case запрос
	useCaseReq := UpdateSubscriptionInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	// Вызов use case
	sub, err := h.subscriptionUseCase.UpdateSubscription(c.Request.Context(), id, useCaseReq)
	if err != nil {
		slog.Error("Failed to update subscription", "id", id, "error", err)
		handleError(c, err)
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

	// Вызов use case
	err := h.subscriptionUseCase.DeleteSubscription(c.Request.Context(), id)
	if err != nil {
		slog.Error("Failed to delete subscription", "id", id, "error", err)
		handleError(c, err)
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

	// Преобразование HTTP запроса в use case запрос
	useCaseReq := ListFiltersInput{
		UserID:      userID,
		ServiceName: serviceName,
		Limit:       limit,
		Offset:      offset,
	}

	// Вызов use case
	subs, err := h.subscriptionUseCase.ListSubscriptions(c.Request.Context(), useCaseReq)
	if err != nil {
		slog.Error("Failed to list subscriptions", "error", err)
		handleError(c, err)
		return
	}

	// Преобразование доменных моделей в HTTP ответы
	responses := make([]SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		responses = append(responses, ToSubscriptionResponse(sub))
	}

	slog.Info("subscriptions listed", "count", len(responses))
	RespondSuccess(c, http.StatusOK, responses)
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

	// Преобразование HTTP запроса в use case запрос
	useCaseReq := SummaryFiltersInput{
		UserID:      userID,
		ServiceName: serviceName,
		PeriodStart: periodStartQuery,
		PeriodEnd:   periodEndQuery,
	}

	// Вызов use case
	total, err := h.subscriptionUseCase.GetSubscriptionsSummary(c.Request.Context(), useCaseReq)
	if err != nil {
		slog.Error("Failed to get summary", "error", err)
		handleError(c, err)
		return
	}

	slog.Info("summary calculated",
		"total", total,
		"from", periodStartQuery,
		"to", periodEndQuery,
	)

	RespondSuccess(c, http.StatusOK, gin.H{
		"total":     total,
		"from":      periodStartQuery,
		"to":        periodEndQuery,
		"user_id":   userID,
		"service":   serviceName,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// handleError обрабатывает ошибки от use case и возвращает соответствующий HTTP статус
func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	errMsg := err.Error()

	// Определяем тип ошибки по содержимому сообщения
	switch {
	case strings.Contains(errMsg, "not found"):
		RespondError(c, http.StatusNotFound, errMsg)
	case strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "required") || strings.Contains(errMsg, "must be"):
		RespondError(c, http.StatusBadRequest, errMsg)
	default:
		RespondError(c, http.StatusInternalServerError, "internal server error")
	}
}
