# Order Matching Engine – Go + MySQL

A simplified stock exchange style order matching engine built in Go. Supports buy/sell, market/limit orders, full price-time priority, partial fills, and trade logging via a REST API.

## Features

- Limit & Market orders for buy/sell
- Price-Time Priority (FIFO within price level)
- Matching logic: partial fills, resting orders, market remainders
- Cancel orders, view order book and trades
- MySQL persistence (no ORM)
- Built with Go + Gin

## Quick Start

### 1. Clone & Enter Project

```bash
git clone https://github.com/samiraghav/GOLANG-ORDER-MATCHING-SYSTEM.git
```

### 2. Install Go Modules

```bash
go mod tidy
```

### 3. MySQL Setup

Create a database (e.g. order_engine) and set credentials in a `.env` file:

#### .env Example

```
DB_USER=root
DB_PASS=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=order_engine
```

Then run the schema SQL:

```bash
mysql -u root -p order_engine < schema.sql
```

### 4. Run the Server

```bash
go run main.go
```

```
Connected to MySQL
Listening and serving HTTP on :8080
```

## API Endpoints

### Place Order

```bash
curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" -d '{
  "symbol": "ETHUSD",
  "side": "buy",
  "type": "limit",
  "price": 2500,
  "quantity": 2
}'
```

### Cancel Order

```bash
curl -X DELETE http://localhost:8080/orders/1
```

### Get Order Book

```bash
curl http://localhost:8080/orderbook?symbol=ETHUSD
```

### List Trades

```bash
curl http://localhost:8080/trades?symbol=ETHUSD
```

### Get Order Status

```bash
curl http://localhost:8080/orders/1
```

## Design Decisions & Assumptions

- Matching happens async after placing an order
- Price-time priority handled via Go's sort.SliceStable
- Market orders match best price or get canceled if unmatched
- No ORM — raw SQL for full control and performance
- created_at is stored as DATETIME, parsed with parseTime=true
- No user authentication — kept simple for assignment scope

## Tech Stack

- Go (Golang)
- Gin Web Framework
- MySQL
- SQL schema (no ORM)
- Raw SQL Queries + Transactions

## Folder Structure

```
├── api/           # REST endpoint handlers
├── db/            # DB layer (raw SQL)
├── models/        # Data models
├── service/       # Matching engine logic
├── utils/         # Response helpers & logging
├── schema.sql     # DB schema
└── main.go        # Entry point
```

## Contact
Reach out via GitHub Issues
