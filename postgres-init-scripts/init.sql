CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    order_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS barcodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    user_id UUID NOT NULL,
    product_id INT NOT NULL,
    code VARCHAR(100) NOT NULL,
    marked boolean DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS shelves (
    id  SERIAL PRIMARY KEY,
    name VARCHAR(10) UNIQUE NOT NULL,
    user_id UUID,
    order_id UUID,
    current_occupancy INT DEFAULT 0 NOT NULL,
    capacity INT DEFAULT 0 NOT NULL
);

CREATE INDEX ON orders (id);

CREATE INDEX ON order_items (id);

CREATE INDEX ON barcodes (code);

CREATE INDEX ON shelves (name);

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "barcodes" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "shelves" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

INSERT INTO shelves (name) VALUES 
    ('A1'),
    ('A2'),
    ('A3'),
    ('A4'),
    ('B1'),
    ('B2'),
    ('B3'),
    ('B4'),
    ('C1'),
    ('C2'),
    ('C3'),
    ('C4');

-- INSERT INTO products (name, description, price) VALUES 
--     ('Bird s Nest Fern', 'The Bird s Nest Fern is a tropical plant known for its vibrant green, wavy fronds...',22),
--     ('Ctenanthe', 'The Ctenanthe, also known as the Prayer Plant, is a stunning tropical plant with bold...',45);

