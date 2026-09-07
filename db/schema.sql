CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY,
    country VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    zip VARCHAR(12) NOT NULL    
)

CREATE TABLE IF NOT EXISTS orders(
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    address_id UUID NOT NULL REFERENCES addresses(id),
    status VARCHAR(20) NOT NULL,
    total BIGINT NOT NULL
)

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    quantity BIGINT NOT NULL,
    price BIGINT NOT NULL    
)

CREATE TABLE IF NOT EXISTS order_outbox (
    id UUID PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    published_at TIMESTAMPTZ
)
