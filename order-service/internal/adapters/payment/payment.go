package payment

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lamtrinh/go-ecom-hexagon/order-service/config"
	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"
	"github.com/sony/gobreaker"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	"github.com/lamtrinh/ecom-proto/go/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
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

	tlsCredentials, tlsErr := getTLSCredentials()
	if tlsErr != nil {
		log.Fatalf("failed to get tls credentials, err: %v", tlsErr)
	}

	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))

	conn, err := grpc.Dial(connection, opts...)
	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentClient(conn)

	return &Adapter{
		payment: client,
	}, nil
}

func getTLSCredentials() (credentials.TransportCredentials, error) {
	certDir := config.GetCertDir()

	cert, certErr := tls.LoadX509KeyPair(certDir+"/order-cert.pem", certDir+"/order-key.pem")
	if certErr != nil {
		return nil, fmt.Errorf("failed to load cert")
	}

	certPool := x509.NewCertPool()
	caCert, caCertErr := os.ReadFile(certDir + "/ca-cert.pem")

	if caCertErr != nil {
		return nil, fmt.Errorf("failed to read ca cert")
	}

	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to append ca cert")
	}

	return credentials.NewTLS(&tls.Config{
		ServerName:   "*.microservices.dev",
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}), nil
}

func (a Adapter) Charge(order *domain.Order) error {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
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
