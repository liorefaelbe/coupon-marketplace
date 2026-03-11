# Digital Coupon Marketplace

Backend service for managing and selling digital coupons to resellers.

The system allows administrators to create and manage digital coupons,
while resellers can browse available products and purchase coupons
through an authenticated API.

---

TECH STACK

- Go
- Gin (HTTP framework)
- PostgreSQL
- Docker
- REST API

---

ARCHITECTURE

Client
↓
HTTP API (Gin)
↓
Service Layer
↓
Repository Layer
↓
PostgreSQL

---

PROJECT STRUCTURE

backend
├─ cmd
│ └─ server
│ └─ main.go
│
├─ internal
│ ├─ handlers
│ ├─ services
│ ├─ repository
│ ├─ middleware
│ └─ database
│
├─ db
│ └─ migrations
│
├─ go.mod
└─ go.sum

---

RUNNING THE PROJECT

1. Clone the repository

git clone https://github.com/liorefaelbe/coupon-marketplace.git
cd coupon-marketplace

2. Start the services

docker compose up --build

3. The API will be available at

http://localhost:8080

---

ENVIRONMENT VARIABLES

The reseller API requires a Bearer token.

Example token used in development:

RESELLER_TOKEN=super-secret-reseller-token

Example request header:

Authorization: Bearer super-secret-reseller-token

---

HEALTH CHECK

Endpoint:

GET /health

Example response:

{
"status": "ok"
}

---

ADMIN API

Create coupon

POST /admin/coupons

Example body:

{
"name": "Amazon $100",
"description": "Gift card",
"image_url": "https://image.com",
"cost_price": 80,
"margin_percentage": 25,
"value_type": "STRING",
"value": "ABCD-1234"
}

List products

GET /admin/products

Get product

GET /admin/products/{id}

Update product

PUT /admin/products/{id}

Delete product

DELETE /admin/products/{id}

---

RESELLER API

Requires header:

Authorization: Bearer super-secret-reseller-token

List available products

GET /api/v1/products

Get product details

GET /api/v1/products/{id}

Purchase coupon

POST /api/v1/products/{id}/purchase

Example request body:

{
"reseller_price": 120
}

Example response:

{
"product_id": "uuid",
"final_price": 120,
"value_type": "STRING",
"value": "ABCD-1234"
}

---

BUSINESS LOGIC

The system enforces the following rules:

- Coupons cannot be sold twice
- Reseller price must be greater than or equal to the minimum selling price
- Coupon purchase uses a database transaction to prevent race conditions

---

DATABASE

PostgreSQL is used for persistence.

Tables:

products
coupons

Each coupon is linked to a product.

---

PROJECT PURPOSE

This project demonstrates:

- layered backend architecture
- transactional database operations
- REST API design
- containerized development environment
