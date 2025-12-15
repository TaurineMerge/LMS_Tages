# Документация: База данных knowledge-base-db

## Обзор

Контейнер `knowledge-base-db` представляет собой базу данных PostgreSQL, предназначенную для хранения контента образовательной платформы (категории, курсы и уроки).

## Конфигурация (Переменные окружения)

Для конфигурации сервиса и управления доступом используются переменные окружения, которые должны быть определены в файле `.env` в корне проекта.

| Переменная | Описание | Пример |
|---|---|---|
| `KNOWLEDGE_BASE_DB_NAME` | Имя создаваемой базы данных. | `knowledge_base` |
| `KNOWLEDGE_BASE_DB_SUPERUSER` | Имя суперпользователя для инициализации БД. | `superuser` |
| `KNOWLEDGE_BASE_DB_SUPERUSER_PASSWORD` | Пароль суперпользователя. | `supersecret` |
| `KNOWLEDGE_BASE_ADMIN_USER` | Имя пользователя-администратора схемы `knowledge_base`. | `kb_admin` |
| `KNOWLEDGE_BASE_ADMIN_PASSWORD` | Пароль пользователя-администратора. | `adminsecret` |
| `KNOWLEDGE_BASE_RO_USER` | Имя пользователя только для чтения. | `kb_ro` |
| `KNOWLEDGE_BASE_RO_PASSWORD` | Пароль пользователя только для чтения. | `rosecret` |
| `KNOWLEDGE_BASE_DB_HOST` | Хост базы данных для подключения из других сервисов. | `knowledge-base-db` |
| `KNOWLEDGE_BASE_DB_PORT` | Внутренний порт базы данных. | `5432` |

### Пример файла .env

```ini
KNOWLEDGE_BASE_DB_NAME=knowledge_base
KNOWLEDGE_BASE_DB_SUPERUSER=superuser
KNOWLEDGE_BASE_DB_SUPERUSER_PASSWORD=supersecret
KNOWLEDGE_BASE_ADMIN_USER=kb_admin
KNOWLEDGE_BASE_ADMIN_PASSWORD=adminsecret
KNOWLEDGE_BASE_RO_USER=kb_ro
KNOWLEDGE_BASE_RO_PASSWORD=rosecret
KNOWLEDGE_BASE_DB_HOST=knowledge-base-db
KNOWLEDGE_BASE_DB_PORT=5432
```

## Инициализация и Миграции

При первом запуске контейнера (когда том с данными пуст), выполняются скрипты из директории `init-sql/knowledge-base-db` в алфавитном порядке.

1.  **`01-init-users.sh`**
    *   Создает схему `knowledge_base`.
    *   Создает роли `KNOWLEDGE_BASE_ADMIN_USER` и `KNOWLEDGE_BASE_RO_USER`.
    *   Настраивает права доступа для этих ролей.
    *   Устанавливает **правила по умолчанию** (`ALTER DEFAULT PRIVILEGES`), чтобы новые таблицы, созданные суперпользователем или `kb_admin`, автоматически давали права на чтение для `kb_ro`.

2.  **`02-create-tables.sql`**
    *   Создает структуру таблиц в схеме `knowledge_base`: `category_d`, `course_b`, `lesson_d`.

3.  **`999-populate.sql`**
    *   Заполняет созданные таблицы тестовыми данными для разработки и отладки.

## Пользователи и Права доступа

Скрипт `01-init-users.sh` создает следующих пользователей с определенными правами:

1.  **Суперпользователь (`${KNOWLEDGE_BASE_DB_SUPERUSER}`)**
    *   **Права:** Полные права на все в базе данных.
    *   **Назначение:** Используется для первоначальной инициализации, создания схемы, таблиц и других ролей.

2.  **Администратор (`${KNOWLEDGE_BASE_ADMIN_USER}`)**
    *   **Права:**
        *   `CONNECT`: Может подключаться к базе данных.
        *   `USAGE ON SCHEMA knowledge_base`: Может "видеть" схему `knowledge_base` и ее объекты.
        *   `SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA knowledge_base`: Может выполнять операции CRUD (чтение, создание, обновление, удаление данных) на всех таблицах в схеме `knowledge_base`.
    *   **Назначение:** Используется сервисом `admin-panel` для управления данными (контентом) в пределах схемы `knowledge_base`. **Не имеет прав на изменение структуры таблиц (DDL операции).**

3.  **Пользователь только для чтения (`${KNOWLEDGE_BASE_RO_USER}`)**
    *   **Права:**
        *   `CONNECT`: Может подключаться к базе данных.
        *   `USAGE`: Может "видеть" схему `knowledge_base` и ее объекты.
        *   `SELECT`: Может выполнять `SELECT` запросы ко всем таблицам в схеме `knowledge_base`.
    *   **Назначение:** Используется сервисом `public-side` для безопасного чтения данных без возможности их изменить.

## Схема данных

Основная логика данных находится в схеме `knowledge_base`. ER-диаграмма, описывающая таблицы и их связи, находится в файле [erd.md](erd.md).

## Сетевой доступ

*   **Внутренний порт:** Внутри сети Docker сервис доступен по порту `5432`.
*   **Внешний порт:** В `docker-compose-dev.yml` порт `5432` контейнера проброшен на порт `3500` хост-машины для удобства подключения из локальных инструментов (например DBeaver).
