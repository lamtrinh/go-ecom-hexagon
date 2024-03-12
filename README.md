## Start the database
```bash
docker compose up -d
```

## Run the order service
```bash
DATABASE_URL=root:root@tcp(localhost:3306)/ecom \
PAYMENT_URL=localhost:3001 \
PORT=3000 \
ENV=development \
go run cmd/main.go
```

## Test Order/Create using grpcurl
```bash
grpcurl -d '{"user_id": 123, "order_items": [{"product_code": "product-123", "quantity": 2, "unit_price": 4}]}' -plaintext localhost:3000  Order/Create
```