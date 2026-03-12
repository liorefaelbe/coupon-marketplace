# Digital Coupon Marketplace

A full-stack digital coupon marketplace built with Go, PostgreSQL, React and Docker.

The system supports three flows:

- Admin can create and manage coupon products
- Direct Customer can purchase coupons directly from the website
- Resellers can integrate through a REST API

---

TECH STACK

Backend

- Go
- Gin
- PostgreSQL

Frontend

- React
- TypeScript
- Vite

Infrastructure

- Docker
- Docker Compose

---

FEATURES

Admin

- Create coupons
- View products
- Update products
- Delete unsold products
- Protected with authentication

Direct Customer

- View available coupons
- Purchase coupons directly from the website
- Receive coupon value only after successful purchase

Reseller API

- Get available products
- Get product by ID
- Purchase coupon with reseller price validation
- Protected with API token authentication

---

PRICING RULES

Each coupon has:

cost_price
margin_percentage

The backend calculates:

minimum_sell_price = cost_price \* (1 + margin_percentage / 100)

Rules enforced by the backend:

- cost_price >= 0
- margin_percentage >= 0
- minimum_sell_price is never accepted from client input
- reseller_price must be >= minimum_sell_price

---

## Project Structure
```text
coupon-marketplace
├── backend
│   ├── cmd
│   │   └── server
│   │       └── main.go
│   ├── db
│   │   └── migrations
│   │       ├── 001_init.sql
│   │       └── 002_seed_data.sql
│   └── internal
│       ├── database
│       │   └── db.go
│       ├── handlers
│       ├── middleware
│       ├── models
│       ├── repository
│       └── services
│
├── frontend
│   ├── public
│   └── src
│       ├── api
│       ├── assets
│       ├── components
│       ├── pages
│       └── types
│
├── docker-compose.yml
├── README.md
└── .gitignore
```
---

RUNNING THE PROJECT

From the project root run:

docker compose up --build

This will start:

- PostgreSQL database
- Go backend API
- React frontend

---

APPLICATION URLS

Frontend
http://localhost:5173

Backend API
http://localhost:8080

Health Check
http://localhost:8080/health

---

DEFAULT CREDENTIALS

Admin login

username: admin
password: admin123

Reseller API token

super-secret-reseller-token

---

API OVERVIEW

ADMIN API

Create coupon
POST /admin/coupons

Get all products
GET /admin/products

Get product by ID
GET /admin/products/{id}

Update product
PUT /admin/products/{id}

Delete product
DELETE /admin/products/{id}

---

CUSTOMER API

Get available products
GET /store/products

Get product by ID
GET /store/products/{id}

Purchase coupon
POST /store/products/{id}/purchase

---

RESELLER API

Get available products
GET /api/v1/products

Get product by ID
GET /api/v1/products/{id}

Purchase product
POST /api/v1/products/{id}/purchase

Header required

Authorization: Bearer super-secret-reseller-token

Example body

{
"reseller_price": 120
}

---

EXAMPLE ADMIN COUPON CREATION

POST /admin/coupons

{
"name": "Amazon $100",
"description": "Gift card",
"image_url": "https://example.com/image.jpg",
"cost_price": 80,
"margin_percentage": 25,
"value_type": "STRING",
"value": "ABCD-1234"
}

---

EXAMPLE PURCHASE RESPONSE

{
"product_id": "uuid",
"final_price": 100,
"value_type": "STRING",
"value": "ABCD-1234"
}

---

SEED DATA

Database tables and seed data are created automatically when running:

docker compose down -v
docker compose up --build

---

NOTES

- Coupon value is returned only after successful purchase
- Direct customers always buy at the displayed price
- Resellers cannot sell below the minimum allowed selling price
- Backend follows layered architecture:
  handlers → services → repository → database

---

Author
Lior Refael Berkovits
