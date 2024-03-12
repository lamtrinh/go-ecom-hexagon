package ports

import (
	"context"

	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/application/domain"
)

type DBPort interface {
	Get(context.Context, int32) (domain.Payment, error)
	Save(context.Context, *domain.Payment) error
}
