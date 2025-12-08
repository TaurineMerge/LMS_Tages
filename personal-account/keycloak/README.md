# Keycloak Configuration

## Автоматическая настройка Realm

При первом запуске Keycloak необходимо импортировать realm-конфигурацию.

### Вариант 1: Импорт через Admin Console

1. Откройте админ-консоль Keycloak: http://localhost:8080
2. Войдите с учетными данными: `admin` / `admin`
3. Нажмите "Create Realm"
4. Выберите "Import" и загрузите файл `student-realm.json`

### Вариант 2: Импорт через командную строку

```bash
# Скопируйте файл в контейнер
docker cp keycloak/student-realm.json keycloak:/tmp/student-realm.json

# Импортируйте realm
docker exec keycloak /opt/keycloak/bin/kc.sh import --file /tmp/student-realm.json
```

### Вариант 3: Автоимпорт при запуске (docker-compose)

Добавьте в `docker-compose.yml` для сервиса keycloak:

```yaml
volumes:
  - ./personal-account/keycloak:/opt/keycloak/data/import
command:
  - start-dev
  - --http-port=8080
  - --import-realm
```

## Конфигурация Realm

### Realm: `student`

- **Регистрация пользователей**: включена
- **Email как username**: выключено
- **Подтверждение email**: выключено (для разработки)

### Роли

| Роль | Описание |
|------|----------|
| `student` | Базовая роль для студентов |
| `teacher` | Роль для преподавателей (полный доступ) |

### Клиенты

#### personal-account-client

- **Client ID**: `personal-account-client`
- **Client Secret**: `personal-account-secret`
- **Redirect URIs**: 
  - `http://localhost/account/api/v1/auth/callback`
  - `http://localhost:8000/account/api/v1/auth/callback`
- **Web Origins**: `*`
- **Protocol**: OpenID Connect
- **Flow**: Authorization Code

### Тестовые пользователи

| Username | Email | Password | Роли |
|----------|-------|----------|------|
| `student1` | student1@example.com | `student123` | student |
| `teacher1` | teacher1@example.com | `teacher123` | teacher, student |

## API Endpoints

### Авторизация

| Method | Endpoint | Описание |
|--------|----------|----------|
| GET | `/account/api/v1/auth/login` | Редирект на страницу входа Keycloak |
| GET | `/account/api/v1/auth/register` | Редирект на страницу регистрации Keycloak |
| POST | `/account/api/v1/auth/register` | Регистрация через API |
| GET | `/account/api/v1/auth/callback` | Обработка callback после авторизации |
| POST | `/account/api/v1/auth/refresh` | Обновление access token |
| POST | `/account/api/v1/auth/logout` | Выход (инвалидация refresh token) |
| GET | `/account/api/v1/auth/me` | Информация о текущем пользователе |

### Регистрация через API

```bash
curl -X POST http://localhost/account/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "securepassword",
    "first_name": "New",
    "last_name": "User"
  }'
```

### Логин через OAuth2 Flow

1. Перейдите на `/account/api/v1/auth/login`
2. Авторизуйтесь в Keycloak
3. Получите токен в callback `/account/api/v1/auth/callback?code=...`

## Environment Variables

```env
# Keycloak Configuration
KEYCLOAK_SERVER_URL=http://keycloak:8080
KEYCLOAK_PUBLIC_URL=http://localhost:8080
KEYCLOAK_REALM=student
KEYCLOAK_CLIENT_ID=personal-account-client
KEYCLOAK_CLIENT_SECRET=personal-account-secret
KEYCLOAK_REDIRECT_URI=http://localhost/account/callback
KEYCLOAK_ADMIN_USERNAME=admin
KEYCLOAK_ADMIN_PASSWORD=admin
```

## Страницы Personal Account

| URL | Описание |
|-----|----------|
| `/account/` | Дашборд (главная страница) |
| `/account/login` | Страница входа |
| `/account/register` | Форма регистрации |
| `/account/callback` | Обработка OAuth callback |
| `/account/profile` | Профиль пользователя |
| `/account/certificates` | Сертификаты |
| `/account/visits` | История посещений |
