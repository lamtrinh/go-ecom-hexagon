package main

import (
	"log"

	"github.com/lamtrinh/go-ecom-hexagon/order-service/config"
	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/adapters/db"
	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/adapters/grpc"
	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/adapters/payment"
	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDatabaseURL())
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	paymentAdapter, err := payment.NewAdapter(config.GetPaymentURL())
	if err != nil {
		log.Fatalf("failed to initialize payment stub, err: %v", err)
	}
	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
