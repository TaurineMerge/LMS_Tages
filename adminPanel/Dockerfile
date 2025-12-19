# Этап сборки
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный этап
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем swagger.json
COPY --from=builder /app/docs ./docs/

# Копируем templates и static для SSR
COPY --from=builder /app/templates ./templates/
COPY --from=builder /app/static ./static/

# Все параметры должны передаваться через переменные окружения
# Нет hardcoded значений - безопасность и гибкость
# Обязательные переменные будут проверены при старте приложения

EXPOSE 4000

CMD ["./main"]