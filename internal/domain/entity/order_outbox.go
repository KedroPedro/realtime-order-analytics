package entity

import (
	"time"
	"uuid"
)

type OrderOutbox struct {
	Id          uuid.UUID
	EventType   string
	Payload     []byte
	Status      string
	CreatedAt   time.Time
	PublishedAt time.Time
}
