CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    customer_id UUID NOT NULL
);

CREATE INDEX ON orders (id);

-- INSERT INTO products (name, description, price) VALUES 
--     ('Bird s Nest Fern', 'The Bird s Nest Fern is a tropical plant known for its vibrant green, wavy fronds...',22),
--     ('Ctenanthe', 'The Ctenanthe, also known as the Prayer Plant, is a stunning tropical plant with bold...',45);

