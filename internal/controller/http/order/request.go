package order

import (
	"time"
	"uuid"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
)

type CreateOrderRequest struct {
	OwnerID string           `json:"owner_id"`
	Items   []OrderItemInput `json:"items"`
	Address *AddressInput    `json:"address,omitempty"`
}

type OrderItemInput struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
	Price    int64  `json:"price"`
}

type AddressInput struct {
	Country string `json:"country"`
	City    string `json:"city"`
	ZIP     string `json:"zip"`
}

const (
	orderNewStatus = "new"
)

func (r *CreateOrderRequest) ToEntity() (*entity.Order, error) {
	orderID := uuid.NewV7()
	ownerID, err := uuid.Parse(r.OwnerID)
	if err != nil {
		return nil, err
	}

	var address *entity.Address
	if r.Address != nil {
		address = &entity.Address{
			Id:      uuid.New(),
			Country: r.Address.Country,
			City:    r.Address.City,
			ZIP:     r.Address.ZIP,
		}
	}

	items := make([]entity.OrderItem, 0, len(r.Items))
	for _, it := range r.Items {
		items = append(items, entity.OrderItem{
			Id:       uuid.New(),
			OrderId:  orderID,
			Name:     it.Name,
			Quantity: it.Quantity,
			Price:    it.Price,
		})
	}

	return &entity.Order{
		Id:        orderID,
		OwnerId:   ownerID,
		CreatedAt: time.Now(),
		Items:     items,
		Address:   address,
		Status:    orderNewStatus,
		Total:     calcTotal(items),
	}, nil
}

func calcTotal(items []entity.OrderItem) int64 {
	var sum int64
	for _, it := range items {
		sum += it.Price * it.Quantity
	}
	return sum
}
