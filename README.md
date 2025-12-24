# LMS_Tages

## Observability stack

В составе `docker-compose` теперь доступны сервисы для трассировки запросов всех API:

- **Jaeger** (`jaegertracing/all-in-one`) — UI и хранилище, открыто на `http://localhost:16686`.
- **OpenTelemetry Collector** — общая точка входа (`0.0.0.0:4317`/gRPC и `0.0.0.0:4318`/HTTP) для всех API (Go, Java, Python).

Collector принимает OTLP-трафик и проксирует его в Jaeger. Чтобы задействовать трассировки, сервисам достаточно отправлять данные на `http://otel-collector:4317` внутри общей сети `app-network`.

Запуск локально (из корня проекта):

```bash
docker compose up jaeger otel-collector <service-name>
```

Например, для Python-сервиса `personal-account` трассировки уже включены и будут отображаться в Jaeger под именем `personal-account-api`.

## Базы данных

Проект использует архитектурный паттерн **"База данных на сервис"**, реализованный в рамках одного инстанса PostgreSQL. Контейнер `app-db` выступает в роли хоста для нескольких логически изолированных баз данных, каждая из которых предназначена для своего модуля.

### Конфигурация (`.env`)

Настройка инстанса `app-db` и создаваемых в нем баз данных осуществляется через переменные окружения в корневом файле `.env`.

| Переменная | Описание |
|---|---|
| `APP_DB_SUPERUSER` | Имя главного суперпользователя для инстанса `app-db`. |
| `APP_DB_SUPERUSER_PASSWORD` | Пароль главного суперпользователя. |
| `KNOWLEDGE_BASE_DB_NAME` | Имя БД для сервисов `admin-panel` и `public-side`. |
| `KNOWLEDGE_BASE_DB_USER` | Имя пользователя-владельца для `KNOWLEDGE_BASE_DB_NAME`. |
| `KNOWLEDGE_BASE_DB_PASSWORD`| Пароль для `KNOWLEDGE_BASE_DB_USER`. |
| `KNOWLEDGE_BASE_ADMIN_USER`| Имя пользователя с правами CRUD в схеме `knowledge_base`. |
| `KNOWLEDGE_BASE_ADMIN_PASSWORD`| Пароль для `KNOWLEDGE_BASE_ADMIN_USER`. |
| `KNOWLEDGE_BASE_RO_USER`| Имя пользователя только для чтения в схеме `knowledge_base`. |
| `KNOWLEDGE_BASE_RO_PASSWORD`| Пароль для `KNOWLEDGE_BASE_RO_USER`. |
| `TESTING_DB_NAME` | Имя БД для сервиса `testing`. |
| `TESTING_DB_USER` | Имя пользователя-владельца для `TESTING_DB_NAME`. |
| `TESTING_DB_PASSWORD` | Пароль для `TESTING_DB_USER`. |
| `PERSONAL_ACCOUNT_DB_NAME` | Имя БД для сервиса `personal-account`. |
| `PERSONAL_ACCOUNT_DB_USER`| Имя пользователя-владельца для `PERSONAL_ACCOUNT_DB_NAME`. |
| `PERSONAL_ACCOUNT_DB_PASSWORD`| Пароль для `PERSONAL_ACCOUNT_DB_USER`. |

### Процесс инициализации

При первом запуске контейнера `app-db` (когда его том данных пуст), автоматически выполняет скрипты из директории `init-sql/`.

1.  **`00-init-dbs.sh`**: Этот скрипт выполняется первым. Он читает переменные из `.env` и создает отдельные базы данных (`knowledge_base`, `testing` и `personal_account`) и их пользователей-владельцев.

2.  **`01-run-all-sub-scripts.sh`**: Этот скрипт-оркестратор запускается следом и выполняет миграции для каждой базы данных, заходя в соответствующие поддиректории внутри `init-sql/`.

    **Важно:** Названия директорий (`knowledge-base-db`, `testing-db`, `personal-account-db`) **жестко прописаны** в скрипте `01-run-all-sub-scripts.sh` и не должны меняться.

3.  **Порядок выполнения скриптов внутри поддиректорий** (`init-sql/<db-name>/`):
    *   **Фаза 1: `.sh` скрипты.** Сначала в алфавитном порядке выполняются все `.sh` файлы. Они отвечают за создание низкоуровневых ролей с гранулированными правами (например, `kb_admin` и `kb_ro`).
    *   **Фаза 2: `.sql` скрипты.** Затем в алфавитном порядке выполняются все `.sql` файлы. Они запускаются от имени пользователя-владельца базы данных (например, `KNOWLEDGE_BASE_DB_USER`) и отвечают за создание таблиц, индексов и наполнение данными.
