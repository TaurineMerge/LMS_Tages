# Personal Account - Technical Documentation

Добро пожаловать в техническую документацию сервиса **Personal Account** системы онлайн-образования LMS Tages.

## О сервисе

Personal Account — это FastAPI-приложение, предоставляющее REST API для управления личным кабинетом студента:

- **Управление студентами** — CRUD операции
- **Сертификаты** — выдача и проверка сертификатов
- **Посещаемость** — отслеживание посещений
- **Авторизация** — интеграция с Keycloak OAuth2

## Быстрый старт

### Swagger UI (интерактивная документация API)

Интерактивная документация API доступна по адресу:

- [Swagger UI](/account/docs-swagger) — для тестирования API endpoints

### Структура проекта

```
personal-account/
├── app/
│   ├── core/           # Безопасность, JWT
│   ├── routers/        # API endpoints
│   ├── services/       # Бизнес-логика
│   ├── repositories/   # Слой данных
│   ├── schemas/        # Pydantic модели
│   ├── config.py       # Конфигурация
│   └── database.py     # Подключение к БД
├── main.py             # Точка входа FastAPI
└── mkdocs.yml          # Конфигурация документации
```

## Навигация

- [Архитектура](architecture.md) — обзор архитектуры приложения
- [API Reference](api/index.md) — справочник по всем модулям

## Технологии

| Технология | Назначение |
|------------|------------|
| FastAPI | Web framework |
| Pydantic | Валидация данных |
| PostgreSQL | База данных |
| Keycloak | Авторизация OAuth2/OIDC |
| OpenTelemetry | Трассировка и метрики |
