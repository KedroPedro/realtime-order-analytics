package ordersoutboxrepo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxCleaner struct {
	conn *pgxpool.Pool
}

func NewOutboxCleaner(conn *pgxpool.Pool) *OutboxCleaner {
	return &OutboxCleaner{
		conn: conn,
	}
}

func (oc *OutboxCleaner) Clean(ctx context.Context) error {
	cleanQuery :=
		`
			UPDATE order_outbox
			SET 
				locked_at = NULL,
				locked_by = NULL,
				status = 'pending'
			WHERE status = 'processing'
				AND locked_at < NOW() - INTERVAL '5 minutes'
		`

	_, err := oc.conn.Exec(ctx, cleanQuery)
	return err
}
