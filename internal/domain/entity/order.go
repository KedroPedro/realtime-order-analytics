package entity

import (
	"time"
	"uuid"
)

type Order struct {
	Id        uuid.UUID
	OwnerId   uuid.UUID
	CreatedAt time.Time
	Items     []OrderItem
	Address   *Address
	Total     int64
}

type OrderItem struct {
	Id       uuid.UUID
	OrderId  uuid.UUID
	Name     string
	Quantity int64
	Price    int64
}

type Address struct {
	Id      uuid.UUID
	Country string
	City    string
	ZIP     string
}
