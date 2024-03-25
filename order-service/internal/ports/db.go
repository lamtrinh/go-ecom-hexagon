package ports

import "github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"

//go:generate mockery --name DBPort
type DBPort interface {
	Get(int32) (domain.Order, error)
	Save(*domain.Order) error
}
