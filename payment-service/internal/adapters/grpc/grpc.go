package grpc

import (
	"context"
	"fmt"

	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/application/domain"

	"github.com/lamtrinh/ecom-proto/go/payment"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Adapter) Create(ctx context.Context, request *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	log.WithContext(ctx).Info("creating payment...")

	newPayment := domain.NewPayment(request.UserId, request.OrderId, request.TotalPrice)
	result, err := a.api.Charge(ctx, newPayment)

	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to charge, error: %v", err)).Err()
	}

	return &payment.CreatePaymentResponse{
		PaymentId: result.ID,
	}, nil
}
