package main

import (
	"log"

	"github.com/lamtrinh/go-ecom-hexagon/payment-service/config"
	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/adapters/grpc"
	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/adapters/grpc/db"
	"github.com/lamtrinh/go-ecom-hexagon/payment-service/internal/application/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database, err: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
