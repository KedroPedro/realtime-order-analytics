package postgresql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envPostgresqlConnString = "POSTGRES_CONN_STRING"
)

type Postgres struct {
	connPool *pgxpool.Pool
}

func NewPostgres() (*Postgres, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	connStr := os.Getenv(envPostgresqlConnString)
	if connStr == "" {
		return nil, fmt.Errorf("environment variable %q is not set", envPostgresqlConnString)
	}

	connPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		connPool: connPool,
	}, nil
}

func (p *Postgres) NewOrdersRepo() *OrdersRepo {
	return NewOrdersRepo(p.connPool)
}
