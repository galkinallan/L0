CREATE TABLE IF NOT EXISTS orders
(
    order_uid  varchar PRIMARY KEY NOT NULL,
    order_info jsonb               NOT NULL
);