package dto

type CreateFollowUpRequest struct {
	CustomerId string `json:"customer_id" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=payment_remind feedback_reply promotion"`
	Note       string `json:"note" validate:"required"`
}

type UpdateFollowUpStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending done"`
}