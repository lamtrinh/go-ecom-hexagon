package grpc

import (
	"context"

	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"

	"github.com/lamtrinh/ecom-proto/go/order"
)

func (a Adapter) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var orderItems []domain.OrderItem
	for _, item := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}

	newOrder := domain.NewOrder(request.UserId, orderItems)
	createdOrder, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		OrderId: createdOrder.ID,
	}, nil
}
