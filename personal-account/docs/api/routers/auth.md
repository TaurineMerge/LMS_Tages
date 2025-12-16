# Auth Router

API endpoints для авторизации.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/token` | Получить JWT токен |
| POST | `/api/v1/auth/logout` | Выйти из системы |
| GET | `/api/v1/auth/me` | Информация о текущем пользователе |

## API Reference

::: app.routers.auth
    options:
      show_root_heading: false
      members_order: source
