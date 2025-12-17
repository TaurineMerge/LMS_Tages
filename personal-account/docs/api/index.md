# API Reference

Справочник по всем модулям Personal Account API.

## Структура

### Routers (Endpoints)

HTTP endpoints приложения:

- [Students](routers/students.md) — управление студентами
- [Certificates](routers/certificates.md) — сертификаты
- [Visits](routers/visits.md) — посещаемость
- [Auth](routers/auth.md) — авторизация
- [Health](routers/health.md) — health checks

### Services (Business Logic)

Бизнес-логика:

- [Student Service](services/student.md)
- [Certificate Service](services/certificate.md)
- [Visit Service](services/visit.md)
- [Auth Service](services/auth.md)
- [Keycloak Service](services/keycloak.md)

### Core (Infrastructure)

Инфраструктурные компоненты:

- [Security](core/security.md) — JWT валидация
- [JWT](core/jwt.md) — работа с токенами

### Schemas (Data Models)

Pydantic модели:

- [Student](schemas/student.md)
- [Certificate](schemas/certificate.md)
- [Visit](schemas/visit.md)
- [Common](schemas/common.md) — общие схемы

### Repositories (Data Access)

Слой доступа к данным:

- [Base Repository](repositories/base.md)
- [Student Repository](repositories/student.md)

### Other

- [Database](database.md) — подключение к БД
- [Config](config.md) — конфигурация
- [Telemetry](telemetry.md) — трассировка
