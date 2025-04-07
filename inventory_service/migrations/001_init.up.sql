CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    product_description TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK(price > 0)
);
