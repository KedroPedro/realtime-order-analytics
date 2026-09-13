package clickhouse

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/KedroPedro/realtime-order-analytics/internal/domain/interfaces"
	orderanalytic "github.com/KedroPedro/realtime-order-analytics/internal/infrastructure/clickhouse/order"
)

type Clickhouse struct {
	conn driver.Conn
}

const (
	envClickhouseAddr       = "CLICKHOUSE_ADDR"
	envClickhouseDatabase   = "CLICKHOUSE_DATABASE"
	envClickhouseUsername   = "CLICKHOUSE_USERNAME"
	envClickhousePassword   = "CLICKHOUSE_PASSWORD"
	envClickhouseClientName = "CLICKHOUSE_CLIENT_NAME"
	envClickhouseClientVer  = "CLICKHOUSE_CLIENT_VER"
)

func NewClickhouse() (*Clickhouse, error) {

	opt, err := setupClickhouseOptions()

	conn, err := clickhouse.Open(opt)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Clickhouse{
		conn: conn,
	}, nil
}

func (ch *Clickhouse) NewOrderAnalytic() interfaces.OrderAnalytic {
	return orderanalytic.NewOrderRepository(ch.conn)
}

func setupClickhouseOptions() (*clickhouse.Options, error) {
	addr := strings.TrimSpace(os.Getenv(envClickhouseAddr))
	if addr == "" {
		return nil, fmt.Errorf("environment variable %q is not set", envClickhouseAddr)
	}

	auth, err := setupClickhouseAuth()
	if err != nil {
		return nil, err
	}

	clientInfo, err := setupClickhouseClientInfo()
	if err != nil {
		return nil, err
	}

	return &clickhouse.Options{
		Addr:       []string{addr},
		Auth:       auth,
		ClientInfo: clientInfo,
	}, nil
}

func setupClickhouseClientInfo() (clickhouse.ClientInfo, error) {
	name := strings.TrimSpace(os.Getenv(envClickhouseClientName))
	if name == "" {
		return clickhouse.ClientInfo{}, fmt.Errorf("environment variable %q is not set", envClickhouseClientName)
	}

	ver := strings.TrimSpace(os.Getenv(envClickhouseClientVer))
	if ver == "" {
		return clickhouse.ClientInfo{}, fmt.Errorf("environment variable %q is not set", envClickhouseClientVer)
	}

	return clickhouse.ClientInfo{
		Products: []struct {
			Name    string
			Version string
		}{
			{Name: name, Version: ver},
		}}, nil
}

func setupClickhouseAuth() (clickhouse.Auth, error) {
	username := strings.TrimSpace(os.Getenv(envClickhouseUsername))
	if username == "" {
		return clickhouse.Auth{}, fmt.Errorf("environment variable %q is not set", envClickhouseUsername)
	}

	database := strings.TrimSpace(os.Getenv(envClickhouseDatabase))
	if database == "" {
		return clickhouse.Auth{}, fmt.Errorf("environment variable %q is not set", envClickhouseDatabase)
	}

	password := strings.TrimSpace(os.Getenv(envClickhousePassword))

	return clickhouse.Auth{
		Database: database,
		Username: username,
		Password: password,
	}, nil
}
