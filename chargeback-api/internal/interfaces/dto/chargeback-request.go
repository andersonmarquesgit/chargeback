package dto

type ChargebackRequest struct {
	UserID        string `json:"user_id" validate:"required"`
	TransactionID string `json:"transaction_id" validate:"required"`
	Reason        string `json:"reason" validate:"required"`
}
