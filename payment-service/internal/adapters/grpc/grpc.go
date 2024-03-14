package grpc

import (
	"context"

	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/application/domain"

	"github.com/lamtrinh/ecom-proto/go/payment"
	log "github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Adapter) Create(ctx context.Context, request *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	log.WithContext(ctx).Info("creating payment...")

	var validationError []*errdetails.BadRequest_FieldViolation
	if request.TotalPrice < 1 {
		validationError = append(validationError, &errdetails.BadRequest_FieldViolation{
			Field:       "TotalPrice",
			Description: "TotalPrice must be greater than 0",
		})
	}
	if len(validationError) > 0 {
		st := status.New(400, "invalid payment request")
		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = validationError
		statusWithDetail, _ := st.WithDetails(badReq)
		return nil, statusWithDetail.Err()
	}

	newPayment := domain.NewPayment(request.UserId, request.OrderId, request.TotalPrice)
	result, err := a.api.Charge(ctx, newPayment)

	if err != nil {
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "application",
			Description: "failed to create payment",
		}
		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		paymentStatus := status.New(codes.InvalidArgument, "payment creation failed")
		statusWithDetail, _ := paymentStatus.WithDetails(badReq)

		return nil, statusWithDetail.Err()
	}

	return &payment.CreatePaymentResponse{
		PaymentId: result.ID,
	}, nil
}
