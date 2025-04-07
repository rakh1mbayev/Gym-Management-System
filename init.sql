CREATE TYPE order_status_enum AS ENUM('payed', 'delivered', 'canceled');
CREATE TYPE role_enum AS ENUM('user','admin');

CREATE TABLE users(
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    role role_enum NOT NULL,
);

CREATE TABLE orders(
    order_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK(price>0),
    status order_status_enum NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INT REFERENCES products(product_id) ON DELETE CASCADE,
    quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
    price_per_item DECIMAL(10, 2) NOT NULL CHECK (price_per_item > 0)
);


CREATE TABLE products(
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    product_description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK(price>0)
);