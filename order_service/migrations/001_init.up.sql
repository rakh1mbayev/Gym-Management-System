CREATE TYPE order_status_enum AS ENUM('payed', 'delivered', 'canceled');

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK(price > 0),
    status order_status_enum NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INT NOT NULL,
    quantity DECIMAL(10,2) NOT NULL CHECK(quantity > 0),
    price_per_item DECIMAL(10,2) NOT NULL CHECK(price_per_item > 0)
);
