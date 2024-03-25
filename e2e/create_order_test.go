package e2e

import (
	"context"
	"log"
	"testing"

	"github.com/lamtrinh/ecom-proto/go/order"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCreateOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}

type CreateOrderTestSuite struct {
	suite.Suite
	compose tc.ComposeStack
}

func (c *CreateOrderTestSuite) SetupSuite() {
	compose, err := tc.NewDockerCompose("./resources/docker-compose.yaml")
	if err != nil {
		log.Fatal("failed to load docker-compose.yaml")
	}
	c.compose = compose

	execError := compose.Up(context.Background(), tc.Wait(true))

	if execError != nil {
		log.Fatalf("failed to execute docker-compose.yaml, err: %v", execError)
	}
}

func (c *CreateOrderTestSuite) Test_Should_Create_Order() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("0.0.0.0:3000", opts...)
	if err != nil {
		log.Fatalf("fail to connect to order service, err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderClient(conn)
	createOrderResponse, errCreateOrder := orderClient.Create(context.Background(), &order.CreateOrderRequest{
		UserId: 123,
		OrderItems: []*order.OrderItem{
			{
				ProductCode: "product-123",
				Quantity:    4,
				UnitPrice:   2,
			},
		},
	})

	c.Equal(errCreateOrder, nil)

	getOrderResponse, errGetOrder := orderClient.Get(context.Background(), &order.GetOrderRequest{OrderId: createOrderResponse.OrderId})
	c.Nil(errGetOrder)
	c.Equal(int32(123), getOrderResponse.UserId)
	orderItem := getOrderResponse.OrderItems[0]
	c.Equal("product-123", orderItem.ProductCode)
	c.Equal(int32(4), orderItem.Quantity)
	c.Equal(float32(2), orderItem.UnitPrice)
}

func (c *CreateOrderTestSuite) TearDownSuite() {
	execError := c.compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal)

	if execError != nil {
		log.Fatalf("failed to shutdown docker-compose.yaml, err: %v", execError)
	}
}
