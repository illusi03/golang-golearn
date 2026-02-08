package model

import "time"

type TransactionModel struct {
	ID          int                      `json:"id"`
	TotalAmount int                      `json:"total_amount"`
	CreatedAt   time.Time                `json:"created_at"`
	Details     []TransactionDetailModel `json:"details"`
}

type TransactionDetailModel struct {
	ID            int     `json:"id"`
	TransactionID int     `json:"transaction_id"`
	ProductID     int     `json:"product_id"`
	ProductName   *string `json:"product_name"`
	Quantity      int     `json:"quantity"`
	Subtotal      int     `json:"subtotal"`
}
