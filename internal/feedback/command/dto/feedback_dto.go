package dto

type CreateFeedbackRequest struct {
	CustomerId string `json:"customer_id" validate:"required"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	Comment    string `json:"comment" validate:"required"`
	Category   string `json:"category" validate:"required,oneof=service payment product branch"`
}