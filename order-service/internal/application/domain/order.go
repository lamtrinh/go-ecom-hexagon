package domain

import "time"

type OrderItem struct {
	ProductCode string  `json:"product_code"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

type Order struct {
	ID         int32       `json:"id"`
	CustomerID int32       `json:"customer_id"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
	CreateAt   int64       `json:"create_at"`
}

func (o Order) TotalPrice() float32 {
	var total float32 = 0
	for _, item := range o.OrderItems {
		total += float32(item.Quantity) * item.UnitPrice
	}

	return total
}

func NewOrder(customerId int32, orderItems []OrderItem) Order {
	return Order{
		CreateAt:   time.Now().Unix(),
		Status:     "pending",
		CustomerID: customerId,
		OrderItems: orderItems,
	}
}
