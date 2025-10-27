package subscriptions

import "time"

type Subscription struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   string     `json:"start_date"` // YYYY-MM
	EndDate     *string    `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateInput struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}
