package dto

import "time"

type FollowUpFilter struct {
	Type   string `query:"type"`
	Status string `query:"status"`
	Search string `query:"search"`
	Branch string `query:"branch"`
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
}

type CustomerSub struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Product string `json:"product"`
	Branch  string `json:"branch"`
}

type FollowUpResponse struct {
	Id         string       `json:"id"`
	CustomerId string       `json:"customer_id"`
	Type       string       `json:"type"`
	Note       string       `json:"note"`
	Status     string       `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	Customer   *CustomerSub `json:"customer,omitempty"`
}

type FollowUpListResponse struct {
	Items      []*FollowUpResponse `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}
