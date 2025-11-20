package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/asgard-born/rest_service_subscriptions/pkg/domain"
	"github.com/asgard-born/rest_service_subscriptions/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscriptionUseCase_CreateSubscription(t *testing.T) {
	mockRepo := &mocks.SubscriptionRepository{}
	useCase := NewSubscriptionUseCase(mockRepo)

	tests := []struct {
		name        string
		input       CreateSubscriptionInput
		mockSetup   func()
		expected    *domain.Subscription
		expectedErr error
	}{
		{
			name: "successful creation",
			input: CreateSubscriptionInput{
				ServiceName: "Netflix",
				Price:       1000,
				UserID:      "user-123",
				StartDate:   "2024-01-01",
				EndDate:     "2024-12-31",
			},
			mockSetup: func() {
				expectedSub := &domain.Subscription{
					ID:          "sub-123",
					ServiceName: "Netflix",
					Price:       1000,
					UserID:      "user-123",
					StartDate:   mustParseDate("2024-01-01"),
					EndDate: sql.NullTime{
						Time:  mustParseDate("2024-12-31"),
						Valid: true,
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Subscription")).
					Return(expectedSub, nil)
			},
			expected: &domain.Subscription{
				ID:          "sub-123",
				ServiceName: "Netflix",
				Price:       1000,
				UserID:      "user-123",
				StartDate:   mustParseDate("2024-01-01"),
				EndDate: sql.NullTime{
					Time:  mustParseDate("2024-12-31"),
					Valid: true,
				},
			},
			expectedErr: nil,
		},
		{
			name: "invalid price",
			input: CreateSubscriptionInput{
				ServiceName: "Netflix",
				Price:       -100, // Отрицательная цена
				UserID:      "user-123",
				StartDate:   "2024-01-01",
			},
			mockSetup:   func() {},
			expected:    nil,
			expectedErr: errors.New("price must be positive"),
		},
		{
			name: "missing required fields",
			input: CreateSubscriptionInput{
				ServiceName: "", // Пустое имя сервиса
				Price:       1000,
				UserID:      "",
				StartDate:   "2024-01-01",
			},
			mockSetup:   func() {},
			expected:    nil,
			expectedErr: errors.New("service name is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := useCase.CreateSubscription(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.ServiceName, result.ServiceName)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSubscriptionUseCase_GetSubscription(t *testing.T) {
	mockRepo := &mocks.SubscriptionRepository{}
	useCase := NewSubscriptionUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		expectedSub := &domain.Subscription{
			ID:          "sub-123",
			ServiceName: "Netflix",
			Price:       1000,
			UserID:      "user-123",
		}

		mockRepo.On("GetByID", mock.Anything, "sub-123").Return(expectedSub, nil)

		result, err := useCase.GetSubscription(context.Background(), "sub-123")

		assert.NoError(t, err)
		assert.Equal(t, expectedSub, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, errors.New("subscription not found"))

		result, err := useCase.GetSubscription(context.Background(), "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSubscriptionUseCase_UpdateSubscription(t *testing.T) {
	mockRepo := &mocks.SubscriptionRepository{}
	useCase := NewSubscriptionUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		input := UpdateSubscriptionInput{
			ServiceName: "Netflix Premium",
			Price:       1500,
			StartDate:   "2024-01-01",
			EndDate:     "2024-12-31",
		}

		expectedSub := &domain.Subscription{
			ID:          "sub-123",
			ServiceName: "Netflix Premium",
			Price:       1500,
			UserID:      "user-123",
		}

		mockRepo.On("Update", mock.Anything, "sub-123", mock.AnythingOfType("*domain.Subscription")).
			Return(expectedSub, nil)

		result, err := useCase.UpdateSubscription(context.Background(), "sub-123", input)

		assert.NoError(t, err)
		assert.Equal(t, expectedSub, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSubscriptionUseCase_ListSubscriptions(t *testing.T) {
	mockRepo := &mocks.SubscriptionRepository{}
	useCase := NewSubscriptionUseCase(mockRepo)

	t.Run("with filters", func(t *testing.T) {
		filters := ListFiltersInput{
			UserID:      "user-123",
			ServiceName: "Netflix",
			Limit:       10,
			Offset:      0,
		}

		expectedSubs := []*domain.Subscription{
			{ID: "sub-1", ServiceName: "Netflix", UserID: "user-123"},
			{ID: "sub-2", ServiceName: "Netflix", UserID: "user-123"},
		}

		mockRepo.On("List", mock.Anything, mock.AnythingOfType("domain.ListFilters")).
			Return(expectedSubs, nil)

		result, err := useCase.ListSubscriptions(context.Background(), filters)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})
}

// Вспомогательная функция для парсинга дат
func mustParseDate(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return t
}
