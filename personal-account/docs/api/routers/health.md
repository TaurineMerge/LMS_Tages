# Health Router

Health check endpoints для мониторинга.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Базовый health check |
| GET | `/health/ready` | Readiness probe |
| GET | `/health/live` | Liveness probe |

## API Reference

::: app.routers.health
    options:
      show_root_heading: false
      members_order: source
