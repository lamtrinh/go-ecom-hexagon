package domain

import "time"

type Payment struct {
	ID         int32   `json:"id"`
	CustomerID int32   `json:"customer_id"`
	Status     string  `json:"status"`
	OrderID    int32   `json:"order_id"`
	TotalPrice float32 `json:"total_price"`
	CreatedAt  int64   `json:"created_at"`
}

func NewPayment(customerId int32, orderId int32, totalPrice float32) Payment {
	return Payment{
		CustomerID: customerId,
		Status:     "pending",
		OrderID:    orderId,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now().Unix(),
	}
}
