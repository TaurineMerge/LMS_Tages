# publicSide

Публичный часть для образовательной платформы LMS Tages. Предоставляет доступ к категория, курсам и урокам. Построен на Go с использованием веб-фреймворка Fiber.

## Технологии

-   **Go** (1.25+)
-   **Fiber v2** - Веб-фреймворк
-   **PostgreSQL** - База данных
-   **Docker** - Контейнеризация
-   **OpenTelemetry** - Сбор трассировок
-   **Jaeger** - Система для визуализации трассировок
-   **slog** - Структурированное логирование

## Запуск и разработка

### Требования

-   Go 1.25 или выше
-   Docker и Docker Compose
-   Настроенный и запущенный `docker-compose.yml` из корня проекта, который предоставляет базу данных и Jaeger/OTel-коллектор.

### Конфигурация

Приложение настраивается через переменные окружения. Их можно задать напрямую в системе или определить в файле `.env` в директории `publicSide`.

**Приоритет:** Переменные, установленные в операционной системе, **переопределяют** значения из `.env` файла.

#### Переменные окружения

| Переменная                    | Описание                                                                            | Обязательная          | Значение по умолчанию |
| ----------------------------- | ----------------------------------------------------------------------------------- | --------------------- | --------------------- |
| `DATABASE_URL`                | Полная строка подключения к PostgreSQL.                                             | Да (или группа `DB_*`) | -                     |
| `DB_HOST`                     | Хост базы данных.                                                                   | Да (если нет `DB_URL`)  | -                     |
| `DB_PORT`                     | Порт базы данных.                                                                   | Да (если нет `DB_URL`)  | -                     |
| `DB_USER`                     | Имя пользователя для подключения к БД.                                              | Да (если нет `DB_URL`)  | -                     |
| `DB_PASSWORD`                 | Пароль для подключения к БД.                                                        | Да (если нет `DB_URL`)  | -                     |
| `DB_NAME`                     | Имя базы данных.                                                                    | Да (если нет `DB_URL`)  | -                     |
| `DB_SSLMODE`                  | Режим SSL для подключения к БД.                                                     | Нет                   | `disable`             |
| `APP_PORT`                    | Порт, на котором будет запущен веб-сервер.                                          | Нет                   | `3000`                |
| `LOG_LEVEL`                   | Уровень логирования (`DEBUG`, `INFO`, `WARN`, `ERROR`).                               | Нет                   | `INFO`                |
| `CORS_ALLOWED_ORIGINS`        | Разрешенные источники для CORS (через запятую).                                     | Нет                   | `*`                   |
| `CORS_ALLOWED_METHODS`        | Разрешенные методы для CORS (через запятую).                                        | Нет                   | `GET` |
| `CORS_ALLOWED_HEADERS`        | Разрешенные заголовки для CORS (через запятую).                                     | Нет                   | `Origin,Content-Type,Accept,Authorization` |
| `CORS_ALLOW_CREDENTIALS`      | Разрешает передачу credentials в CORS.                                              | Нет                   | `false`                |
| `OTEL_SERVICE_NAME`           | Имя сервиса, которое будет отображаться в Jaeger.                                   | Нет                   | `publicSide`          |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Адрес OTel-коллектора для отправки трассировок (например, `localhost:4317`).        | Да                    | -                     |

#### Пример `.env` файла

Создайте файл `publicSide/.env` по этому шаблону:

```env
# Database
DB_USER=appuser
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=appdb

# Application
APP_PORT=3000
LOG_LEVEL=DEBUG

# CORS
CORS_ALLOWED_ORIGINS="http://localhost:3000,http://localhost:9090"
CORS_ALLOWED_METHODS="GET,POST,PUT,DELETE,OPTIONS"
CORS_ALLOWED_HEADERS="Origin, Content-Type, Accept, Authorization"
CORS_ALLOW_CREDENTIALS=true

# OpenTelemetry
OTEL_SERVICE_NAME=publicSide-local
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
```

### Запуск через Docker Compose

Это предпочтительный способ запуска приложения вместе с зависимостями (база данных, Jaeger/OTel Collector).

1.  **Запустите все сервисы** (включая `publicSide`):
    ```bash
    # Убедитесь, что вы находитесь в корневой директории проекта (LMS_Tages)
    docker compose up --build -d public-side
    ```
    *Эта команда соберет образ `publicSide` (если он изменился) и запустит его в фоновом режиме.*

2.  **Просмотр логов:**
    ```bash
    docker compose logs -f public-side
    ```

4.  **Остановка сервисов:**
    ```bash
    docker compose down
    ```

### Локальный запуск (без Docker)

1.  **Установите зависимости:**
    ```bash
    # Убедитесь, что вы находитесь в директории publicSide
    go mod tidy
    ```

2.  **Создайте и настройте файл `.env`** (см. секцию "Конфигурация").

3.  **Запустите сервер:**
    ```bash
    go run ./cmd/main.go
    ```

4.  Приложение будет доступно по адресу, указанному в `APP_PORT` (по умолчанию `http://localhost:3000`).

## API

Документация по API доступна в формате Swagger. После запуска сервера перейдите по адресу:
[http://localhost:3000/api/v1/swagger/index.html](http://localhost:3000/api/v1/swagger/index.html)

## Линтинг и качество кода

Проект использует `golangci-lint` для статического анализа и поддержания качества кода. Подробные инструкции по установке и использованию находятся в файле [doc/linter.md](./doc/linter.md).

## Структура проекта

```
publicSide/
├── cmd/main.go         # Точка входа в приложение
├── internal/           # Внутренняя логика, не предназначенная для импорта другими проектами
│   ├── config/
│   ├── handler/
│   ├── repository/
│   └── service/
├── pkg/                # Независимые пакеты
│   ├── apiconst/
│   ├── apperrors/
│   ├── database/
│   └── tracing/
├── doc/                # Документация (Swagger, Linter)
├── .env                # Файл с переменными окружения (не должен быть в git)
├── go.mod
└── Dockerfile
```
