package db

import (
	"context"
	"fmt"

	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/application/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	CustomerID int32
	Status     string
	OrderID    int32
	TotalPrice float32
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(connection string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	mgErr := db.AutoMigrate(&Payment{})
	if mgErr != nil {
		return nil, fmt.Errorf("db migration error: %v", mgErr)
	}

	return &Adapter{db: db}, nil
}

func (a Adapter) Get(ctx context.Context, id int32) (domain.Payment, error) {
	var paymentModel Payment
	res := a.db.First(&paymentModel, id)
	payment := domain.Payment{
		ID:         int32(paymentModel.ID),
		CustomerID: paymentModel.CustomerID,
		Status:     paymentModel.Status,
		OrderID:    paymentModel.OrderID,
		TotalPrice: paymentModel.TotalPrice,
		CreatedAt:  paymentModel.CreatedAt.UnixMicro(),
	}

	return payment, res.Error
}

func (a Adapter) Save(ctx context.Context, payment *domain.Payment) error {
	paymentModel := Payment{
		CustomerID: payment.CustomerID,
		Status:     payment.Status,
		OrderID:    payment.OrderID,
		TotalPrice: payment.TotalPrice,
	}

	res := a.db.WithContext(ctx).Create(&paymentModel)
	if res.Error == nil {
		payment.ID = int32(paymentModel.ID)
	}
	return res.Error
}
