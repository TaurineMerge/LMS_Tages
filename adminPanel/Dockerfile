# Этап сборки
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Копируем ВСЕ файлы сразу
COPY . .

# Синхронизируем и собираем
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный этап
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 4000

CMD ["./main"]