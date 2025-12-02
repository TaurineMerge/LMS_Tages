# publicSide

Публичная часть LMS Tages, построенная на Go Fiber.

## Требования

- Docker

## Как запустить Docker с сервером

### Запуск через Docker

1. **Соберите Docker образ:**
```bash
docker build -t publicside-server .
```

2. **Запустите контейнер:**
```bash
docker run -p 3000:3000 publicside-server
```

3. **Откройте браузер:**
```
http://localhost:3000
```

### Использование другого порта

Если порт 3000 занят, можно использовать другой:
```bash
docker run -p 8080:3000 publicside-server
```
Тогда приложение будет доступно на `http://localhost:8080`

## Разработка без Docker

1. **Установите зависимости:**
```bash
go mod tidy
```

2. **Запустите сервер:**
```bash
go run main.go
```

3. **Приложение будет доступно на:**
```
http://localhost:3000
```

## Структура проекта

```
publicSide/
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## Технологии

- [Go](https://golang.org/) - язык программирования
- [Fiber v3](https://gofiber.io/) - веб-фреймворк
- [Docker](https://www.docker.com/) - контейнеризация