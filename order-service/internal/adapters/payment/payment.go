package payment

import (
	"context"
	"time"

	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"

	"github.com/lamtrinh/ecom-proto/go/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(connection string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(connection, opts...)
	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentClient(conn)

	return &Adapter{
		payment: client,
	}, nil
}

func (a Adapter) Charge(order *domain.Order) error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	return err
}
