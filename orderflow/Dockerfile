FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /order-api ./cmd/order-api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /order-api /app/order-api

EXPOSE 8080

CMD ["/app/order-api"]
