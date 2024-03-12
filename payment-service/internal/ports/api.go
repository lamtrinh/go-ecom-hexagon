package ports

import (
	"context"

	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/application/domain"
)

type APIPort interface {
	Charge(context.Context, domain.Payment) (domain.Payment, error)
}
