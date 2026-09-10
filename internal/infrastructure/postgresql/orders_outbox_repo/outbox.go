package ordersoutboxrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdersOutboxRepo struct {
	conn *pgxpool.Pool
}

func New(connPool *pgxpool.Pool) *OrdersOutboxRepo {
	return &OrdersOutboxRepo{
		conn: connPool,
	}
}

const (
	batchSize = 10
)

func (or *OrdersOutboxRepo) GetEventsBatch(ctx context.Context, locker string) ([]entity.OrderOutbox, error) {
	tx, err := or.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	getOrders := fmt.Sprintf(
		`
			SELECT id, event_type, payload, created_at FROM order_outbox
			WHERE published_at IS NULL AND locked_at IS NULL
			ORDER BY created_at ASC
			LIMIT %d
			FOR UPDATE SKIP LOCKED
		`, batchSize)

	rows, err := tx.Query(
		ctx,
		getOrders,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	temp := entity.OrderOutbox{}

	batch := make([]entity.OrderOutbox, 0, batchSize)

	for rows.Next() {
		err := rows.Scan(
			&temp.Id,
			&temp.EventType,
			&temp.Payload,
			&temp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		batch = append(batch, temp)
	}

	lockOrders :=
		`
			UPDATE order_oubtox
			SET 
				locked_at = $1,
				locked_by = $2,
				attempts = attempts + 1,
				status = 'processing'
			WHERE id = $3
		`

	updateBatch := &pgx.Batch{}

	for _, ev := range batch {
		updateBatch.Queue(
			lockOrders,
			time.Now(), locker, ev.Id,
		)
	}

	res := tx.SendBatch(ctx, updateBatch)
	for range updateBatch.Len() {
		if _, err := res.Exec(); err != nil {
			return nil, err
		}
	}

	if err := res.Close(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return batch, nil
}

func (or *OrdersOutboxRepo) PublishEvents(ctx context.Context, events []entity.OrderOutbox) error {
	tx, err := or.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	publishEvents :=
		`
			UPDATE order_outbox
			SET published_at = $1
			WHERE id = $2
		`

	batch := pgx.Batch{}
	for _, event := range events {
		batch.Queue(
			publishEvents,
			time.Now, event.Id,
		)
	}

	res := tx.SendBatch(ctx, &batch)

	for range batch.Len() {
		if _, err := res.Exec(); err != nil {
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
