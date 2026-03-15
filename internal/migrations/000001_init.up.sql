CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(20) NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE cars (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    vin TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX idx_cars_user_id ON cars(user_id);
CREATE UNIQUE INDEX idx_user_car ON cars(user_id);

CREATE TABLE cart_items (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    part_id VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    brand VARCHAR(50) NOT NULL,
    price BIGINT NOT NULL,
    quantity BIGINT NOT NULL,
    delivery_day INT NOT NULL,
    image_url TEXT,
    PRIMARY KEY (user_id, part_id)
);
CREATE INDEX idx_cart_items_user_id ON cart_items(user_id);

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    address VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX idx_orders_user_id ON orders(user_id);

CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    part_id VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    brand VARCHAR(50) NOT NULL,
    price BIGINT NOT NULL,
    quantity INT NOT NULL,
    delivery_day INT NOT NULL
);