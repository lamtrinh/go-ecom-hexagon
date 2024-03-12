package ports

import "github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"

type PaymentPort interface {
	Charge(*domain.Order) error
}