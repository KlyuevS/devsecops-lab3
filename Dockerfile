# -----------------------------
# 1. Сборочный этап
# -----------------------------
FROM golang:1.24-alpine AS builder

# Включаем модули и кэш
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем всё остальное
COPY . .

# Собираем бинарник
RUN go build -o news-app ./cmd/server


# 2. Финальный минимальный образ

FROM alpine:3.19

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/news-app .

# Приложение слушает 8080
EXPOSE 8080

CMD ["./news-app"]
