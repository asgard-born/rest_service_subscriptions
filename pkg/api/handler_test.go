package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asgard-born/rest_service_subscriptions/pkg/domain"
	"github.com/asgard-born/rest_service_subscriptions/pkg/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubscriptionUseCase реализует SubscriptionUseCase для тестов
type MockSubscriptionUseCase struct {
	mock.Mock
}

func (m *MockSubscriptionUseCase) CreateSubscription(ctx context.Context, req usecase.CreateSubscriptionInput) (*domain.Subscription, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionUseCase) GetSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionUseCase) UpdateSubscription(ctx context.Context, id string, req usecase.UpdateSubscriptionInput) (*domain.Subscription, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionUseCase) DeleteSubscription(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionUseCase) ListSubscriptions(ctx context.Context, filters usecase.ListFiltersInput) ([]*domain.Subscription, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionUseCase) GetSubscriptionsSummary(ctx context.Context, filters usecase.SummaryFiltersInput) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func TestHandler_CreateSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockSubscriptionUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success",
			requestBody: CreateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       1000,
				UserID:      "user-123",
				StartDate:   "2024-01-01",
				EndDate:     "2024-12-31",
			},
			mockSetup: func(muc *MockSubscriptionUseCase) {
				muc.On("CreateSubscription", mock.Anything, usecase.CreateSubscriptionInput{
					ServiceName: "Netflix",
					Price:       1000,
					UserID:      "user-123",
					StartDate:   "2024-01-01",
					EndDate:     "2024-12-31",
				}).Return(&domain.Subscription{
					ID:          "sub-123",
					ServiceName: "Netflix",
					Price:       1000,
					UserID:      "user-123",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"id":           "sub-123",
					"service_name": "Netflix",
					"price":        float64(1000),
					"user_id":      "user-123",
					"currency":     "",
				},
			},
		},
		{
			name:           "invalid json",
			requestBody:    `invalid json`,
			mockSetup:      func(muc *MockSubscriptionUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"success": false,
				"error":   "invalid request body",
			},
		},
		{
			name: "validation error - missing required fields",
			requestBody: map[string]interface{}{
				"price": 1000,
				// missing service_name and user_id
			},
			mockSetup:      func(muc *MockSubscriptionUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"success": false,
				"error":   "invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &MockSubscriptionUseCase{}
			tt.mockSetup(mockUC)

			handler := NewHandler(mockUC)

			// Создаем тестовый запрос
			var bodyBytes []byte
			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/subscriptions", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router := gin.New()
			router.POST("/subscriptions", handler.CreateSubscription)
			router.ServeHTTP(rr, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Проверяем тело ответа
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)

			if err != nil {
				return
			}

			if tt.expectedBody["success"].(bool) {
				assert.True(t, response["success"].(bool))
				// Дополнительные проверки данных
			} else {
				assert.False(t, response["success"].(bool))
				assert.Contains(t, response["error"].(string), tt.expectedBody["error"].(string))
			}

			mockUC.AssertExpectations(t)
		})
	}
}

func TestHandler_GetSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &MockSubscriptionUseCase{}
	handler := NewHandler(mockUC)

	t.Run("success", func(t *testing.T) {
		mockUC.On("GetSubscription", mock.Anything, "sub-123").
			Return(&domain.Subscription{
				ID:          "sub-123",
				ServiceName: "Netflix",
				Price:       1000,
				UserID:      "user-123",
			}, nil)

		req := httptest.NewRequest("GET", "/subscriptions/sub-123", nil)
		rr := httptest.NewRecorder()

		router := gin.New()
		router.GET("/subscriptions/:id", handler.GetSubscription)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			return
		}

		assert.True(t, response["success"].(bool))
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "sub-123", data["id"])
	})

	t.Run("not found", func(t *testing.T) {
		mockUC.On("GetSubscription", mock.Anything, "nonexistent").
			Return(nil, errors.New("subscription not found"))

		req := httptest.NewRequest("GET", "/subscriptions/nonexistent", nil)
		rr := httptest.NewRecorder()

		router := gin.New()
		router.GET("/subscriptions/:id", handler.GetSubscription)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestHandler_ListSubscriptions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := &MockSubscriptionUseCase{}
	handler := NewHandler(mockUC)

	t.Run("success with filters", func(t *testing.T) {
		mockUC.On("ListSubscriptions", mock.Anything, usecase.ListFiltersInput{
			UserID:      "user-123",
			ServiceName: "Netflix",
			Limit:       10,
			Offset:      0,
		}).Return([]*domain.Subscription{
			{ID: "sub-1", ServiceName: "Netflix", UserID: "user-123"},
			{ID: "sub-2", ServiceName: "Netflix", UserID: "user-123"},
		}, nil)

		req := httptest.NewRequest("GET", "/subscriptions?user_id=user-123&service_name=Netflix&limit=10&offset=0", nil)
		rr := httptest.NewRecorder()

		router := gin.New()
		router.GET("/subscriptions", handler.ListSubscriptions)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			return
		}

		assert.True(t, response["success"].(bool))
		data := response["data"].([]interface{})
		assert.Len(t, data, 2)
	})

	t.Run("invalid limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/subscriptions?limit=invalid", nil)
		rr := httptest.NewRecorder()

		router := gin.New()
		router.GET("/subscriptions", handler.ListSubscriptions)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{"not found", errors.New("subscription not found"), http.StatusNotFound},
		{"invalid input", errors.New("invalid input data"), http.StatusBadRequest},
		{"internal error", errors.New("internal server error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			rr := httptest.NewRecorder()

			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				handleError(c, tt.err)
			})
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
