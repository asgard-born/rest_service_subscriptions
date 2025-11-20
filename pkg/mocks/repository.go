package mocks

import (
	"context"

	"github.com/asgard-born/rest_service_subscriptions/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type SubscriptionRepository struct {
	mock.Mock
}

func (m *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	args := m.Called(ctx, sub)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *SubscriptionRepository) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *SubscriptionRepository) Update(ctx context.Context, id string, sub *domain.Subscription) (*domain.Subscription, error) {
	args := m.Called(ctx, id, sub)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *SubscriptionRepository) List(ctx context.Context, filters domain.ListFilters) ([]*domain.Subscription, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Subscription), args.Error(1)
}

func (m *SubscriptionRepository) GetSummary(ctx context.Context, filters domain.SummaryFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}
