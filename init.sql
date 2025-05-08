CREATE EXTENSION IF NOT EXISTS pgcrypto;


CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50),
    confirmation_token VARCHAR(255),
    is_confirmed BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    product_description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price_per_item DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (order_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (product_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_order_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_item_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_item_product_id ON order_items(product_id);

INSERT INTO users (name, email, password, role, is_confirmed)
VALUES
    ('Petr', 'petr@gmail.com', '$2a$10$ZlC9n6kDeY/rF7BkOIx1..sdyGaJqdrQ0F2NhWk2iU.1oZ6XB/U2a', 'user', true),
    ('Pavel', 'pavel@gmail.com', '$2a$10$LfKx7R4dCVv2ZrbGZQUyoehZkEkcGkPbG9.9IzG0xZfq.zjCv6KQe', 'admin', true);

-- "password123" | "admin123".

INSERT INTO products (name, product_description, price, stock)
VALUES
    ('Protein Powder', 'Whey protein powder for muscle building', 39.99, 100),
    ('Yoga Mat', 'Non-slip yoga mat for workouts', 25.00, 50),
    ('Dumbbell Set', 'Adjustable dumbbells up to 40kg', 89.99, 30);

INSERT INTO orders (user_id, total_price, status)
VALUES (1, 64.99, 'pending')
RETURNING order_id;

