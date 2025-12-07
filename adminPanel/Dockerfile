# Этап сборки
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Устанавливаем переменные окружения по умолчанию
ENV ROOT_PATH=/admin
ENV API_PREFIX=/api/v1
ENV API_ADDRESS=:4000
ENV CORS_ALLOW_ORIGINS=http://localhost,http://localhost:3000
ENV CORS_ALLOW_CREDENTIALS=false

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный этап
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем swagger.json
COPY --from=builder /app/docs ./docs/

# Устанавливаем переменные окружения
ENV ROOT_PATH=/admin
ENV API_PREFIX=/api/v1
ENV API_ADDRESS=:4000
ENV CORS_ALLOW_ORIGINS=http://localhost,http://localhost:4000
ENV CORS_ALLOW_CREDENTIALS=false

EXPOSE 4000

CMD ["./main"]