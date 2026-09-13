SET allow_experimental_json_type = 1;

CREATE DATABASE IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.kafka_order_events(
    event JSON
) ENGINE = Kafka()
SETTINGS
    kafka_broker_list = 'kafka:9092',
    kafka_topic_list = 'order_outbox',
    kafka_group_name = 'clickhouse_group',
    kafka_format = 'JSONAsObject',
    kafka_num_consumers = 3; 

CREATE TABLE IF NOT EXISTS analytics.order_events (
    id UUID,
    event_type LowCardinality(String),
    created_at DateTime64(3),
    order_id UUID,
    owner_id UUID,
    total Decimal(10,2)
) ENGINE = MergeTree
ORDER BY (id, event_type, created_at);

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.events_kafka_mv
TO analytics.order_events
AS
SELECT
    toUUID(event.id) as id,
    parseDateTime64BestEffort(
        event.created_at::String,
        3
    ) as created_at,
    event.event_type::String as event_type,
    toUUID(event.payload.order_id) as order_id,
    toUUID(event.payload.owner_id) as owner_id,
    toDecimal128(event.payload.total::String, 2) as total
FROM analytics.kafka_order_events;
    
