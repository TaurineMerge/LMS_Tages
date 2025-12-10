# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Копируем файлы модуля
COPY go.mod go.sum ./

# Явно скачиваем зависимости
RUN go get github.com/gofiber/fiber/v2@v2.52.0 && \
    go get github.com/gofiber/template/handlebars/v2@v2.1.7

# Копируем исходный код
COPY . .

# Синхронизируем зависимости и генерируем go.sum
RUN go mod tidy

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем ВСЕ файлы (включая шаблоны в правильной структуре)
COPY --from=builder /app/layouts ./layouts
COPY --from=builder /app/pages ./pages

EXPOSE 3000

CMD ["./main"]