package model

import "time"

type OnlineBusiness struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	WebsiteName string    `json:"website_name"`
	URL         string    `json:"url"`
	EnamadID    string    `json:"enamad_id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
