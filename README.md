## Start the database
```
docker compose up -d
```

## Run the order service
```
DATABASE_URL=root:root@tcp(localhost:3306)/ecom \
PAYMENT_URL=localhost:3001 \
PORT=3000 \
ENV=development \
CERT_DIR=[generated certificates directory] \
go run cmd/main.go
```

## Run the payment service
```
DATABASE_URL=root:root@tcp(localhost:3306)/ecom \
PORT=3001 \
ENV=development \
CERT_DIR=[generated certificates directory] \
go run cmd/main.go
```

## Test Order/Create using grpcurl
```
grpcurl -d '{"user_id": 123, "order_items": [{"product_code": "product-123", "quantity": 2, "unit_price": 4}]}' -plaintext localhost:3000  Order/Create
```

## mTLS certificate generation
### Generate a private key and a self-signed certificate for the certificate authority (CA)
```
openssl req -x509 \
    -sha256 \
    -newkey rsa:4096 \
    -days 365 \
    -keyout ca-key.pem \
    -out ca-cert.pem \
    -subj "/C=VN/ST=Ho Chi Minh City/L=Ho Chi Minh City/O=Software/OU=Microservices/CN=*.microservices.dev/emailAddress=go@microservices.dev" \
    -nodes
```

#### Verify the generated self-certificate for the CA
```
openssl x509 -in ca-cert.pem -noout -text
```
### Generate PaymentService private key and certificate signing request
```
openssl req \
    -newkey rsa:4096 \
    -keyout payment-key.pem \
    -out payment-req.pem \
    -subj "/C=VN/ST=Ho Chi Minh City/L=Ho Chi Minh City/O=Software/OU=PaymentService/CN=*.microservices.dev/emailAddress=go@microservices.dev" \
    -nodes \
    -sha256
```
### Sign it using the CA’s private key
```
openssl x509 \
    -req -in payment-req.pem \
    -days 60 \
    -CA ca-cert.pem \
    -CAkey ca-key.pem \
    -CAcreateserial \
    -out payment-cert.pem \
    -extfile payment-ext.cnf \
    -sha256
```
#### Example configuration for ext file option
```
subjectAltName=DNS:*.microservices.dev,DNS:*.microservices.dev,IP:0.0.0.0
```
#### Verify the PaymentService’s self-signed certificate
```
openssl x509 -in payment-cert.pem -noout -text
```
### Generate OrderService private key and certificate signing request
```
openssl req \
    -newkey rsa:4096 \
    -keyout order-key.pem \
    -out order-req.pem \
    -subj "/C=VN/ST=Ho Chi Minh City/L=Ho Chi Minh City/O=Software/OU=OrderService/CN=*.microservices.dev/emailAddress=go@microservices.dev" \
    -nodes \
    -sha256
```
#### Sign it using the CA’s private key
```
openssl x509 \
    -req -in order-req.pem \
    -sha256 \
    -days 60 \
    -CA ca-cert.pem \
    -CAkey ca-key.pem \
    -CAcreateserial \
    -out order-cert.pem \
    -extfile order-ext.cnf
```
#### Verify the OrderService’s self-signed certificate
```
openssl x509 -in order-cert.pem -noout -text
```