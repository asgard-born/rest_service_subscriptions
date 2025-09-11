package api

import "github.com/asgard-born/rest_service_subscriptions/pkg/db"

func ToSubscriptionResponse(s db.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
