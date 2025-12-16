# Personal Account API

API микросервис для личного кабинета системы онлайн образования.

## Технологии

- **FastAPI** — асинхронный веб-фреймворк
- **Pydantic** — валидация данных
- **Psycopg3** — драйвер PostgreSQL (без ORM, чистый SQL)
- **Uvicorn** — ASGI сервер
- **OpenTelemetry** — сбор трассировки запросов

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

### Authentication
| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/api/v1/auth/login` | Редирект на Keycloak |
| GET | `/api/v1/auth/register` | Редирект на регистрацию Keycloak |
| POST | `/api/v1/auth/register` | Регистрация через API |
| GET | `/api/v1/auth/callback` | Обработка OAuth callback |
| POST | `/api/v1/auth/refresh` | Обновление токена |
| POST | `/api/v1/auth/logout` | Выход |
| GET | `/api/v1/auth/me` | Текущий пользователь |

### Frontend Pages
| URL | Описание |
|-----|----------|
| `/` | Дашборд |
| `/login` | Страница входа |
| `/register` | Форма регистрации |
| `/callback` | OAuth callback |
| `/profile` | Профиль |
| `/certificates` | Сертификаты |
| `/visits` | Посещения |

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
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Endpoint OTLP (gRPC/HTTP) | `http://otel-collector:4317` |
| `OTEL_SERVICE_NAME` | Имя сервиса в трассировках | `personal-account-api` |
| `OTEL_EXPORTER_OTLP_INSECURE` | Отключение TLS для OTLP | `true` |
| `KEYCLOAK_SERVER_URL` | URL Keycloak (внутренний) | `http://keycloak:8080` |
| `KEYCLOAK_PUBLIC_URL` | URL Keycloak (публичный) | `http://localhost:8080` |
| `KEYCLOAK_REALM` | Realm Keycloak | `student` |
| `KEYCLOAK_CLIENT_ID` | Client ID | `personal-account-client` |
| `KEYCLOAK_CLIENT_SECRET` | Client Secret | `personal-account-secret` |
| `KEYCLOAK_REDIRECT_URI` | Redirect URI | `http://localhost/account/callback` |
| `KEYCLOAK_ADMIN_USERNAME` | Admin username | `admin` |
| `KEYCLOAK_ADMIN_PASSWORD` | Admin password | `admin` |

## Настройка Keycloak

Для работы авторизации необходимо импортировать realm в Keycloak:

```bash
# После запуска контейнеров
docker cp personal-account/keycloak/student-realm.json keycloak:/tmp/student-realm.json
docker exec keycloak /opt/keycloak/bin/kc.sh import --file /tmp/student-realm.json
```

Или через Admin Console (http://localhost:8080):
1. Войдите как `admin`/`admin`
2. Create Realm → выберите файл `keycloak/student-realm.json`

Подробнее: [keycloak/README.md](keycloak/README.md)

## Observability

- Трассы автоматически отправляются в OpenTelemetry Collector (см. `observability/otel-collector-config.yaml`).
- По умолчанию контейнер отправляет данные на `http://otel-collector:4317`, откуда они уезжают в Jaeger (`http://localhost:16686`).
- Инструментируются уровни FastAPI → сервисы → репозитории/SQLAlchemy, поэтому в Jaeger виден полный путь запроса вплоть до отдельных SQL выражений.
- Чтобы проверить локально, запускайте из корня репозитория:

```bash
docker compose up personal-account otel-collector jaeger
```

После старта откройте Jaeger UI `http://localhost:16686` и выберите сервис `personal-account-api`.

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
