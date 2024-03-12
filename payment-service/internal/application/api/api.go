package api

import (
	"context"

	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/application/domain"
	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error) {
	err := a.db.Save(ctx,  &payment)
	if err != nil{
		return domain.Payment{}, err
	}

	return payment, nil
}
