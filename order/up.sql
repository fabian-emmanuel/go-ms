CREATE TABLE IF NOT EXISTS orders (
    id CHAR(30) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    account_id CHAR(30) NOT NULL,
    total_amount MONEY NOT NULL
);

CREATE TABLE IF NOT EXISTS ordered_products (
    order_id CHAR(30) REFERENCES orders(id) ON DELETE CASCADE,
    product_id CHAR(30) NOT NULL,
    quantity INT NOT NULL,
    PRIMARY KEY (order_id, product_id)
);