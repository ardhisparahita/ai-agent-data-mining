# Dockerfile untuk ecommerce-api (Go)
# Taruh file ini di root repo Go Anda (sejajar dengan go.mod)

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Sesuaikan port ini dengan port yang dipakai aplikasi Go Anda (cek main.go / config)
EXPOSE 3000

CMD ["./main"]
