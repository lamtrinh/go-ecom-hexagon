package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/lamtrinh/go-ecom-hexagon/order-service/config"
	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/ports"

	"github.com/lamtrinh/ecom-proto/go/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api  ports.APIPort
	port int
	order.UnimplementedOrderServer
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

	server := grpc.NewServer()
	order.RegisterOrderServer(server, a)

	if config.GetEnv() == "development" {
		reflection.Register(server)
	}

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d", a.port)
	}
}
