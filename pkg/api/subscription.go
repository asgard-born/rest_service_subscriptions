package api

import "time"

type SubscriptionResponse struct {
	ID          string     `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       float64    `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
