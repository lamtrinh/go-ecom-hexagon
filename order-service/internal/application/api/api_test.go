package api

import (
	"errors"
	"testing"

	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"
	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Should_Place_Order(t *testing.T) {
	payment := ports.NewMockPaymentPort(t)
	db := ports.NewMockDBPort(t)
	payment.On("Charge", mock.Anything).Return(nil)
	db.On("Save", mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "product-123",
				UnitPrice:   2,
				Quantity:    4,
			},
		},
	})

	assert.Nil(t, err)
}

func Test_Should_Return_Error_When_Db_Fail(t *testing.T) {
	payment := ports.NewMockPaymentPort(t)
	db := ports.NewMockDBPort(t)
	db.On("Save", mock.Anything).Return(errors.New("connection error"))

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "product-123",
				UnitPrice:   2,
				Quantity:    4,
			},
		},
	})

	assert.EqualError(t, err, "connection error")
}

func Test_Should_Return_Error_When_Payment_Fail(t *testing.T) {
	payment := ports.NewMockPaymentPort(t)
	db := ports.NewMockDBPort(t)
	payment.On("Charge", mock.Anything).Return(errors.New("insufficient balance"))
	db.On("Save", mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "product-123",
				UnitPrice:   2,
				Quantity:    4,
			},
		},
	})
	st, _ := status.FromError(err)

	assert.Equal(t, st.Message(), "order creation failed")
	assert.Equal(t, st.Details()[0].(*errdetails.BadRequest).FieldViolations[0].Description, "insufficient balance")
	assert.Equal(t, st.Code(), codes.InvalidArgument)
}
