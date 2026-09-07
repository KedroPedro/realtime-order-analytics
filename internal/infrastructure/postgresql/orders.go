package postgresql

import (
	"context"
	"encoding/json"
	"time"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdersRepo struct {
	conn *pgxpool.Pool
}

func NewOrdersRepo(connPool *pgxpool.Pool) *OrdersRepo {
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

	insertAddress :=
		`
		INSERT INTO addresses (id, country, city, zip)
		VALUES ($1, $2, $3, %4)
		`

	if _, err := tx.Exec(
		ctx,
		insertAddress,
		order.Address.Id, order.Address.Country,
		order.Address.City, order.Address.ZIP,
	); err != nil {
		return err
	}

	insertItem :=
		`
		INSERT INTO order_items (id, order_id, name, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
		`

	batch := &pgx.Batch{}
	for _, item := range order.Items {
		batch.Queue(
			insertItem,
			item.Id, item.OrderId, item.Name, item.Quantity, item.Price,
		)
	}

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

	insertOrder :=
		`
		INSERT INTO orders (id, owner_id, created_at, address_id, status, total)
		VALUES ($1, $2, $3, $4, $5) 
		`

	if _, err := tx.Exec(
		ctx,
		insertOrder,
		order.Id, order.OwnerId, order.CreatedAt,
		order, order.Status, order.Total,
	); err != nil {
		return err
	}

	insertOutbox :=
		`
		INSERT INTO order_outbox (id, event_type, payload, created_at, published_at)
		VALUES ($1, $2, $3, $4, $5)
		`

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

	if _, err := tx.Exec(
		ctx,
		insertOutbox,
		order.Id, order.event_type, payload,
		order.CreatedAt, time.Now,
	); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
