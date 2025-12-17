# Инструкция по работе с Docker

## Структура путей в Docker контейнере

После запуска через Docker Compose, в контейнере `public-side` файлы находятся по следующим путям:

```
/root/
├── main                    # Исполняемый файл
├── views/                  # Handlebars шаблоны
│   ├── layouts/
│   │   └── main.hbs
│   ├── partials/
│   │   └── course-card.hbs
│   └── pages/
│       └── courses.hbs
├── static/                 # Статические файлы
│   └── css/
│       └── styles.css
└── doc/                    # Swagger документация
    └── swagger/
```

## Запуск через Docker Compose

1. **Сборка и запуск всех сервисов:**
```bash
docker-compose up -d --build
```

2. **Просмотр логов publicSide:**
```bash
docker-compose logs -f public-side
```

3. **Проверка статуса контейнеров:**
```bash
docker-compose ps
```

4. **Доступ к приложению через Nginx:**
```
http://localhost/categories/{category_id}/courses
```

Nginx проксирует запросы к `public-side` контейнеру.

## Отладка

### Проверить файлы внутри контейнера:
```bash
docker exec -it public-side ls -la /root/
docker exec -it public-side ls -la /root/views/
docker exec -it public-side ls -la /root/static/
```

### Зайти внутрь контейнера:
```bash
docker exec -it public-side sh
```

### Проверить переменные окружения:
```bash
docker exec -it public-side env
```

## Важные моменты

1. **Пути в коде относительные** - `./views`, `./static`, `./doc`
2. **Рабочая директория** в контейнере - `/root/`
3. **Порт внутри контейнера** - 3000 (настроен в `.env.prod`)
4. **Доступ снаружи** - через Nginx на порту 80

## Переменные окружения

Файл `publicSide/.env.prod` должен содержать:
```env
# Database
DATABASE_HOST=app-db
DATABASE_PORT=5432
DATABASE_USER=appuser
DATABASE_PASSWORD=password
DATABASE_NAME=appdb
DATABASE_POOL_MIN_SIZE=5
DATABASE_POOL_MAX_SIZE=20

# Server
PORT=3000

# CORS
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_ALLOW_CREDENTIALS=false

# OpenTelemetry
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
OTEL_EXPORTER_OTLP_PROTOCOL=grpc
OTEL_SERVICE_NAME=public-side

# Logging
LOG_LEVEL=INFO
```

## Обновление кода

После изменения кода нужно пересобрать контейнер:
```bash
docker-compose up -d --build public-side
```

## Проверка доступности

1. **Проверка health endpoint** (если есть):
```bash
docker exec -it public-side wget -qO- http://localhost:3000/health
```

2. **Проверка через curl внутри контейнера:**
```bash
docker exec -it public-side wget -qO- http://localhost:3000/api/v1/swagger/
```

3. **Проверка через хост:**
```bash
curl http://localhost/api/v1/swagger/
```
