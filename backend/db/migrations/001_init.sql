CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    type TEXT NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE coupons (
    product_id UUID PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    cost_price NUMERIC NOT NULL,
    margin_percentage NUMERIC NOT NULL,
    minimum_sell_price NUMERIC NOT NULL,
    is_sold BOOLEAN DEFAULT FALSE,
    value_type TEXT NOT NULL,
    value TEXT NOT NULL
);