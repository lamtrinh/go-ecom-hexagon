package ports

import "github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"

type APIPort interface {
	PlaceOrder(domain.Order) (domain.Order, error)
	GetOrder(int32) (domain.Order, error)
}
