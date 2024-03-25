package db

import (
	"fmt"

	"github.com/lamtrinh/go-ecom-hexagon/order-service/internal/application/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID int32
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(connection string) (*Adapter, error) {
	db, conErr := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if conErr != nil {
		return nil, fmt.Errorf("db connection error: %v", conErr)
	}

	mgErr := db.AutoMigrate(&Order{}, &OrderItem{})
	if mgErr != nil {
		return nil, fmt.Errorf("db migration error: %v", mgErr)
	}

	return &Adapter{
		db: db,
	}, nil
}

func (a Adapter) Get(id int32) (domain.Order, error) {
	var orderModel Order
	res := a.db.Preload("OrderItems").First(&orderModel, id)
	var orderItems []domain.OrderItem
	for _, orderItem := range orderModel.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	order := domain.Order{
		ID:         int32(orderModel.ID),
		CustomerID: int32(orderModel.CustomerID),
		Status:     orderModel.Status,
		OrderItems: orderItems,
		CreateAt:   orderModel.CreatedAt.UnixNano(),
	}

	return order, res.Error
}

func (a Adapter) Save(order *domain.Order) error {
	var orderItemsModel []OrderItem
	for _, orderItem := range order.OrderItems {
		orderItemsModel = append(orderItemsModel, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	orderModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItemsModel,
	}
	res := a.db.Save(&orderModel)

	if res.Error == nil {
		order.ID = int32(orderModel.ID)
	}

	return res.Error
}
