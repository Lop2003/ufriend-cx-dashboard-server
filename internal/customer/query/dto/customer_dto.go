package dto

import "time"

type CustomerFilter struct {
	Branch    string `query:"branch"`
	Status    string `query:"status"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
}

type CustomerResponse struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Product    string    `json:"product"`
	Branch     string    `json:"branch"`
	PlanMonths int       `json:"plan_months"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type FeedbackSub struct {
	Id        string    `json:"id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Category  string    `json:"category"`
	Sentiment string    `json:"sentiment"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowUpSub struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Note      string    `json:"note"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomerDetailResponse struct {
	Id         string        `json:"id"`
	Name       string        `json:"name"`
	Phone      string        `json:"phone"`
	Product    string        `json:"product"`
	Branch     string        `json:"branch"`
	PlanMonths int           `json:"plan_months"`
	Status     string        `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	Feedbacks  []FeedbackSub `json:"feedbacks"`
	FollowUps  []FollowUpSub `json:"follow_ups"`
}

type SummaryResponse struct {
	TotalCustomers int     `json:"total_customers"`
	AvgRating      float64 `json:"avg_rating"`
	OverdueCount   int     `json:"overdue_count"`
	ActiveCount    int     `json:"active_count"`
	CompletedCount int     `json:"completed_count"`
}

type BranchStat struct {
	Branch        string  `json:"branch" bson:"_id"`
	CustomerCount int     `json:"customer_count" bson:"customer_count"`
	AvgRating     float64 `json:"avg_rating" bson:"avg_rating"`
	OverdueCount  int     `json:"overdue_count" bson:"overdue_count"`
}