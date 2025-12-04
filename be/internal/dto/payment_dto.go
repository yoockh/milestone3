package dto

import "milestone3/be/internal/entity"

type PaymentRequest struct {
	UserId int `json:"user_id"`
	AuctionItemId float64 `json:"auction_item_id"`
	Amount float64 `json:"amount" validate:"required"`
}

type PaymentResponse struct {
	PaymentLinkUrl string `json:"payment_link_url"`
	TransactionId string `json:"transaction_id"`
	ExpiryTime string `json:"expiry_time"`
	OrderId string `json:"order_id"`
}

type CheckPaymentStatusResponse struct {
	OrderId string `json:"order_id"`
	TransactionId string `json:"transaction_id"`
	PaymentStatus string `json:"payment_status"`
}

type PaymentInfoResponse struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	User entity.Users `json:"user"`
	AuctionItemId int `json:"auction_item_id"`
	Status string `json:"payment_status"`
	// PaymentStatus entity.PaymentStatus `json:"payment_status"`
	Amount float64 `json:"amount"`
}