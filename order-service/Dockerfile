FROM golang:1.22.0-alpine AS builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order-service ./cmd/main.go

FROM scratch
COPY --from=builder /usr/src/app/order-service ./order-service
CMD [ "./order-service" ]