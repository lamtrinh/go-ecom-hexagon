package payment

import (
	"context"
	"log"
	"time"

	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"
	"github.com/sony/gobreaker"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/lamtrinh/ecom-proto/go/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(connection string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithCodes(codes.Unavailable),
		grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(3*time.Second)),
	)))
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
	// ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err := a.payment.Create(context.TODO(), &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	return err
}

func CircuitBreakerClientInterceptor(cb *gobreaker.CircuitBreaker) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		_, cbErr := cb.Execute(func() (interface{}, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return nil, err
			}

			return nil, nil
		})

		return cbErr
	}
}

func NewCircuitBreakerClientInterceptor() grpc.DialOption {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "payment breaker",
		MaxRequests: 3,
		Timeout:     4,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.TotalSuccesses)
			return failureRatio >= 0.6
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("circuit breaker: %s, changed from %v to %v", name, from, to)
		},
	})

	return grpc.WithUnaryInterceptor(CircuitBreakerClientInterceptor(cb))
}
