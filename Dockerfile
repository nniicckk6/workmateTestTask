# syntax=docker/dockerfile:1
# --- Стадия сборки ---
FROM golang:1.24-alpine AS builder
WORKDIR /app
# Копируем модули и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download
# Копируем исходники и собираем статический бинарь
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o workmateTestProject main.go

# --- Финальная стадия ---
FROM alpine:3.18
# Создаём рабочую директорию и копируем из builder
WORKDIR /app
COPY --from=builder /app/workmateTestProject /app/workmateTestProject


# По умолчанию эти переменные можно переопределять при запуске контейнера
ENV PORT=8080
ENV MAX_CONCURRENT_TASKS=10


# Открываем порт
EXPOSE 8080


# Запуск сервиса
ENTRYPOINT ["/app/workmateTestProject"]
