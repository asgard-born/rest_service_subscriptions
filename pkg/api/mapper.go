package api

import "github.com/asgard-born/rest_service_subscriptions/pkg/db"

func ToSubscriptionResponse(s db.Subscription) SubscriptionResponse {
	var endDate string
	if s.EndDate.Valid {
		endDate = s.EndDate.Time.Format("01-2006")
	}

	return SubscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   s.StartDate.Format("01-2006"),
		EndDate:     endDate,
		CreatedAt:   s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
