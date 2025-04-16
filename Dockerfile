# syntax=docker/dockerfile:1.4

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o pvz-service ./cmd/server/main.go

# Финальный минимальный образ
FROM alpine:3.19

WORKDIR /app

# Копируем собранный бинарник
COPY --from=builder /app/pvz-service .

# Порт, на котором запускается HTTP/gRPC
EXPOSE 8080
EXPOSE 3000
EXPOSE 9000

CMD ["./pvz-service"]