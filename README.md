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
