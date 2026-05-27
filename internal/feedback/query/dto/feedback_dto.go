package dto

import "time"

type FeedbackFilter struct {
	Category string `query:"category"`
	Rating   int    `query:"rating"`
	Branch   string `query:"branch"`
}

type FeedbackResponse struct {
	Id         string    `json:"id"`
	CustomerId string    `json:"customer_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	Category   string    `json:"category"`
	Sentiment  string    `json:"sentiment"`
	CreatedAt  time.Time `json:"created_at"`
}

type FeedbackStatsResponse struct {
	AvgRating     float64   `json:"avg_rating"`
	PositiveCount int       `json:"positive_count"`
	NeutralCount  int       `json:"neutral_count"`
	NegativeCount int       `json:"negative_count"`
	WeeklyCSAT     []float64 `json:"weekly_csat"`
}