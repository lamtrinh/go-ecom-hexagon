package db

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type OrderDatabaseTestSuite struct {
	suite.Suite
	DatabaseURL string
}

func (o *OrderDatabaseTestSuite) SetupSuite() {
	ctx := context.Background()
	port := "3306/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("root:root@tcp(localhost:%s)/ecom?charset=utf8mb4&parseTime=True&loc=Local", port.Port())
	}
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:5.7",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "ecom",
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "mysql", dbURL).WithStartupTimeout(30 * time.Second),
	}

	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal("failed to start mysql test container")
	}
	endpoint, _ := mysqlContainer.Endpoint(ctx, "")
	o.DatabaseURL = fmt.Sprintf("root:root@tcp(%s)/ecom?charset=utf8mb4&parseTime=True&loc=Local", endpoint)
}

func (o *OrderDatabaseTestSuite) Test_Should_Save_Order() {
	adapter, err := NewAdapter(o.DatabaseURL)
	o.Nil(err)
	saveErr := adapter.Save(&domain.Order{})
	o.Nil(saveErr)
}

func (o *OrderDatabaseTestSuite) Test_Should_Get_Order() {
	adapter, _ := NewAdapter(o.DatabaseURL)
	order := domain.NewOrder(1, []domain.OrderItem{
		{
			ProductCode: "product-123",
			Quantity:    4,
			UnitPrice:   2,
		},
	})

	adapter.Save(&order)
	foundOrder, _ := adapter.Get(order.ID)
	o.Equal(int32(1), foundOrder.CustomerID)
}

func TestOrderDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDatabaseTestSuite))
}
