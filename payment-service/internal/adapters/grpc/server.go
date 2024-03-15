package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/lamtrinh/go-ecom-hexagon/payment-service/config"
	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/ports"

	"github.com/lamtrinh/ecom-proto/go/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *grpc.Server
	payment.UnimplementedPaymentServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}

	tlsCredentials, tlsErr := getTLSCredentials()
	if tlsErr != nil {
		log.Fatalf("failed to get tls credentials, err: %v", tlsErr)
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(tlsCredentials))

	grpcServer := grpc.NewServer(opts...)
	a.server = grpcServer
	payment.RegisterPaymentServer(grpcServer, a)

	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	log.Printf("starting payment service on port %d...", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d", a.port)
	}
}

func (a Adapter) Stop() {
	a.server.Stop()
}

func getTLSCredentials() (credentials.TransportCredentials, error) {
	certDir := config.GetCertDir()

	cert, certErr := tls.LoadX509KeyPair(certDir+"/payment-cert.pem", certDir+"/payment-key.pem")
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
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	}), nil
}
