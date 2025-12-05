# Personal Account API

API микросервис для личного кабинета системы онлайн образования.

## Технологии

- **FastAPI** — асинхронный веб-фреймворк
- **Pydantic** — валидация данных
- **Psycopg3** — драйвер PostgreSQL (без ORM, чистый SQL)
- **Uvicorn** — ASGI сервер

## Архитектура

Проект построен по принципу чистой архитектуры с разделением на слои:

```
app/
├── config.py          # Конфигурация приложения
├── database.py        # Пул соединений с БД
├── exceptions.py      # Кастомные исключения
├── schemas/           # Pydantic схемы (DTO)
│   ├── common.py
│   ├── student.py
│   ├── certificate.py
│   └── visit.py
├── repositories/      # Слой работы с БД (raw SQL)
│   ├── base.py
│   ├── student.py
│   ├── certificate.py
│   └── visit.py
├── services/          # Бизнес-логика
│   ├── student.py
│   ├── certificate.py
│   └── visit.py
└── routers/           # API эндпоинты
    ├── students.py
    ├── certificates.py
    ├── visits.py
    └── health.py
```

## API Endpoints

### Students
| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/v1/students` | Список студентов (пагинация) |
| GET | `/api/v1/students/{id}` | Получить студента |
| POST | `/api/v1/students` | Создать студента |
| PUT | `/api/v1/students/{id}` | Обновить студента |
| DELETE | `/api/v1/students/{id}` | Удалить студента |

### Certificates
| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/v1/certificates` | Список сертификатов |
| GET | `/api/v1/certificates/{id}` | Получить сертификат |
| POST | `/api/v1/certificates` | Создать сертификат |
| DELETE | `/api/v1/certificates/{id}` | Удалить сертификат |

### Visits
| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/v1/visits` | Список посещений |
| GET | `/api/v1/visits/{id}` | Получить посещение |
| POST | `/api/v1/visits` | Зарегистрировать посещение |
| DELETE | `/api/v1/visits/{id}` | Удалить посещение |

### Health
| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/health` | Проверка сервиса |
| GET | `/health/db` | Проверка БД |

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_HOST` | Хост БД | `app-db` |
| `DATABASE_PORT` | Порт БД | `5432` |
| `DATABASE_NAME` | Имя БД | `appdb` |
| `DATABASE_USER` | Пользователь | `appuser` |
| `DATABASE_PASSWORD` | Пароль | `password` |
| `DATABASE_POOL_MIN_SIZE` | Мин. пул | `5` |
| `DATABASE_POOL_MAX_SIZE` | Макс. пул | `20` |

## Запуск

### Локально
```bash
pip install -r requirements.txt
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

### Docker
```bash
docker build -t personal-account .
docker run -p 8000:8000 personal-account
```

### Docker Compose (из корня проекта)
```bash
docker-compose up personal-account
```

## Документация API

- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc
- OpenAPI: http://localhost:8000/openapi.json
