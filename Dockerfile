FROM golang:1.26 AS builder

    WORKDIR /app

    COPY go.mod go.sum ./
    RUN go mod download

    COPY . .
    RUN go build -o expense-api ./cmd/api

FROM debian:stable-slim

    WORKDIR /app

    COPY --from=builder /app/expense-api ./expense-api
    CMD ["./expense-api"]

# docker build -t expense-api:multi .
# docker run --rm -p 8080:8080 -e DB_PATH=/data/expenses.db -v expense-data:/data expense-api:multi
