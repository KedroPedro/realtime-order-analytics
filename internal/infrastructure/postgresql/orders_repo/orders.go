package ordersrepo

import (
	"context"
	"encoding/json"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	OrderCreatedEvent = "order_created"
)

type OrdersRepo struct {
	conn *pgxpool.Pool
}

func New(connPool *pgxpool.Pool) *OrdersRepo {
	return &OrdersRepo{
		conn: connPool,
	}
}

func (or *OrdersRepo) CreateOrder(ctx context.Context, order *entity.Order) error {
	tx, err := or.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	insertAddress :=
		`
		INSERT INTO addresses (id, country, city, zip)
		VALUES ($1, $2, $3, $4)
		`

	batch.Queue(
		insertAddress,
		order.Address.Id, order.Address.Country,
		order.Address.City, order.Address.ZIP,
	)

	insertOrder :=
		`
		INSERT INTO orders (id, owner_id, created_at, address_id, total)
		VALUES ($1, $2, $3, $4, $5) 
		`

	batch.Queue(
		insertOrder,
		order.Id, order.OwnerId, order.CreatedAt,
		order.Address.Id, order.Total,
	)

	insertItem :=
		`
		INSERT INTO items (id, order_id, name, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
		`

	for _, item := range order.Items {
		batch.Queue(
			insertItem,
			item.Id, item.OrderId, item.Name,
			item.Quantity, item.Price,
		)
	}

	payload, err := json.Marshal(
		map[string]any{
			"order_id": order.Id,
			"owner_id": order.OwnerId,
			"total":    order.Total,
		},
	)

	if err != nil {
		return err
	}

	insertOutbox :=
		`
		INSERT INTO order_outbox (id, event_type, payload, created_at, status)
		VALUES ($1, $2, $3, $4, 'pending')
		`

	batch.Queue(
		insertOutbox,
		order.Id, OrderCreatedEvent,
		payload, order.CreatedAt,
	)

	res := tx.SendBatch(ctx, batch)

	for i := 0; i < batch.Len(); i++ {
		if _, err := res.Exec(); err != nil {
			res.Close()
			return err
		}
	}

	if err := res.Close(); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
