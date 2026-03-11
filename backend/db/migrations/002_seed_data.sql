INSERT INTO products (id, name, description, type, image_url, created_at, updated_at)
VALUES
('11111111-1111-1111-1111-111111111111','Amazon $50','Amazon gift card','COUPON','https://image.com',NOW(),NOW()),
('22222222-2222-2222-2222-222222222222','Netflix Subscription','Netflix 1 month','COUPON','https://image.com',NOW(),NOW()),
('33333333-3333-3333-3333-333333333333','Steam Wallet $20','Steam wallet code','COUPON','https://image.com',NOW(),NOW());

INSERT INTO coupons (product_id,cost_price,margin_percentage,minimum_sell_price,is_sold,value_type,value)
VALUES
('11111111-1111-1111-1111-111111111111',30,30,39,false,'STRING','AMZN-AAAA-1111'),
('22222222-2222-2222-2222-222222222222',7,40,9.8,false,'STRING','NFLX-BBBB-2222'),
('33333333-3333-3333-3333-333333333333',12,35,16.2,false,'STRING','STM-CCCC-3333');