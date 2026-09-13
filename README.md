# Realtime Order Analytics

Order processing and analytics system using Go, PostgreSQL, Kafka, and ClickHouse.

## Stack

- Go 1.27.0
- PostgreSQL — persistent storage
- Kafka — message broker
- ClickHouse — analytics

## Services

**order-service** — receives orders, saves to PostgreSQL, writes to outbox table

**worker-service** — reads from outbox, publishes to Kafka

**analytic-service** — consumes Kafka events, stores in ClickHouse, serves pre-defined queries
