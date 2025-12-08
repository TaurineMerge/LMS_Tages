# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Repository overview

- This repository is a multi-service LMS prototype composed of several small services orchestrated via Docker Compose and fronted by Nginx.
- Core services:
  - `adminPanel/` – Go Fiber v2 backend exposing an admin REST API over the LMS domain.
  - `publicSide/` – Go Fiber v3 public-facing HTTP service.
  - `personal-account/` – FastAPI service for personal account functionality.
  - `testing/` – Java Javalin demo service.
- Supporting infrastructure:
  - `docker-compose.yml` defines `nginx`, `keycloak` + `keycloak-db`, `app-db` (PostgreSQL) and the four application services on a shared `app-network`.
  - `nginx/nginx.conf` routes traffic on port 80 to the individual backends.
  - Domain language and concepts are documented in `docs/Glossary.md` and implemented in `init-sql/migrate-001.sql` (schemas `personal_account`, `knowledge_base`, `tests`).

## Running the full environment (Docker Compose)

From the repository root, use Docker Compose to build and start the full stack:

```bash path=null start=null
docker compose up --build
```

(On older setups you may need `docker-compose` instead of `docker compose`.)

This will start:
- `nginx` on host ports `80`/`443` (entrypoint for all HTTP traffic).
- `keycloak` on `8080` backed by `keycloak-db` Postgres.
- `app-db` Postgres on `5432`, seeded from `init-sql/migrate-001.sql`.
- `public-side` (Go, port `3000` in container, exposed as `3000:3000`).
- `admin-panel` (Go, port `4000` in container, exposed as `4000:4000`).
- `personal-account` (FastAPI, port `8000` in container, exposed as `8000:8000`).
- `testing` (Javalin, port `8085` in container, exposed as `8085:8085`).

Key external entrypoints (assuming default compose config):
- Public site: `http://localhost/` (proxied to `publicSide` via Nginx).
- Admin API: `http://localhost/admin/` (proxied to `adminPanel`).
- Personal account API: `http://localhost/account/` (proxied to `personal-account`).
- Testing service: `http://localhost/testing/` (proxied to `testing`).
- Keycloak: `http://localhost/auth/`.

To rebuild only a specific service:

```bash path=null start=null
docker compose build admin-panel
# or
docker compose build public-side
```

## Service-level commands

### adminPanel (Go, port 4000)

Location: `adminPanel/`

- Purpose: Admin REST API over the `knowledge_base` domain (categories, courses, lessons) backed by PostgreSQL.
- The service listens on `:4000` and exposes:
  - JSON API under `/api/v1` (e.g. `/api/v1/categories`, `/api/v1/courses`, `/api/v1/courses/{id}/lessons`).
  - Health check at `/api/v1/health`.
  - Swagger UI at `/swagger/index.html` backed by `/swagger.json` (served from `adminPanel/docs/swagger.json`).
- Database configuration:
  - Uses env var `DATABASE_URL`; if unset, falls back to `postgresql://appuser:password@app-db:5432/appdb?sslmode=disable` (matching `app-db` in `docker-compose.yml`).

Run via Docker (from `adminPanel/README.md`):

```bash path=null start=null
cd adminPanel
docker build -t adminpanel .
docker run -p 4000:4000 adminpanel
```

Local dev without Docker (uses the same module as the Docker image):

```bash path=null start=null
cd adminPanel
go mod tidy
go run main.go
```

When running locally without Compose, ensure `DATABASE_URL` points to a reachable Postgres instance (the default assumes the `app-db` container on the compose network).

### publicSide (Go, port 3000)

Location: `publicSide/`

- Purpose: Public-facing Go Fiber v3 HTTP service (currently a simple "Hello, World!" at `/`).
- In Compose, exposed on `http://localhost:3000` and fronted by Nginx at `http://localhost/`.

Run via Docker (from `publicSide/README.md`):

```bash path=null start=null
cd publicSide
docker build -t publicside-server .
docker run -p 3000:3000 publicside-server
```

Local dev without Docker (from `publicSide/README.md`):

```bash path=null start=null
cd publicSide
go mod tidy
go run main.go
```

### personal-account (FastAPI, port 8000)

Location: `personal-account/`

- Purpose: FastAPI service for the personal account / profile domain (currently a minimal `GET /` endpoint, with richer domain modeled in the database schema).
- In Compose, exposed on `http://localhost:8000` and proxied by Nginx under `/account/`.

Quick deploy via Docker (from `personal-account/README.md`):

```bash path=null start=null
cd personal-account
docker build -t fastapi-app .
docker run -d -p 8000:8000 --name fastapi-container fastapi-app
```

Local dev without Docker (mirrors the container CMD):

```bash path=null start=null
cd personal-account
pip install -r requirements.txt
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

### testing (Java Javalin, port 8085)

Location: `testing/`

- Purpose: Small Javalin-based Java service for experimentation/testing with endpoints like `/` and `/api/hello`.
- In Compose, exposed on `http://localhost:8085` and proxied by Nginx under `/testing/`.

Run via Docker (from `testing/README.md`):

```bash path=null start=null
cd testing
docker build -t javalin-app .
docker run javalin-app
```

The `testing/Dockerfile` starts the app using Maven inside the container:

```bash path=null start=null
mvn clean compile exec:java -Dexec.mainClass=com.example.Main
```

You can use the same Maven command locally from `testing/` if you add the appropriate `exec-maven-plugin` configuration to `pom.xml`.

### Supporting services and routing

- `app-db` (PostgreSQL):
  - Defined in `docker-compose.yml` with database `appdb`, user `appuser`, password `password`.
  - Initialized from `init-sql/migrate-001.sql`, which creates schemas and tables for `personal_account`, `knowledge_base`, and `tests`.
- `keycloak` and `keycloak-db`:
  - Used for authentication/authorization, reachable via Nginx under `/auth/` (proxied to `keycloak:8080`).
- `nginx`:
  - Configured in `nginx/nginx.conf` with upstreams for each service (`public-side`, `admin-panel`, `personal-account`, `testing`, `keycloak`).
  - Key locations:
    - `/auth/` → Keycloak backend.
    - `/admin/` → adminPanel backend.
    - `/account/` → personal-account backend.
    - `/testing/` → testing backend.
    - `/` (fallback) → publicSide backend.
  - Shared proxy settings (timeouts, headers, buffers) are in `nginx/includes/proxy_params.conf`.

## Architecture & domain model

### Domain language and contexts

- The ubiquitous language for the LMS domain is defined in `docs/Glossary.md` (Russian descriptions with English identifiers used in code).
- Major domain contexts:
  - **Personal account (`personal_account` schema)** – students, certificates, visit tracking.
  - **Knowledge base (`knowledge_base` schema)** – categories, courses, lessons and their visibility/level.
  - **Tests (`tests` schema)** – tests, questions, answers and test attempts.
- Naming and status conventions:
  - Status and enum-like fields in the database use English values (e.g. `draft`, `public`, `private`, `hard`, `medium`, `easy`), while UI-facing text may be localized.
  - Glossary documents constraints such as visibility states, level values, and ID formats.

### Database schema (PostgreSQL)

Defined in `init-sql/migrate-001.sql` and applied automatically by the `app-db` service:

- Schema `personal_account`:
  - `student_s` – students (profile info, email, contacts).
  - `certificate_b` – certificates linked to students, courses and test attempts.
  - `visit_students_for_lessons` – student lesson visit records.
- Schema `knowledge_base`:
  - `category_d` – course categories.
  - `course_b` – courses with `level` and `visibility` constraints.
  - `lesson_d` – lessons with JSONB `content`.
- Schema `tests`:
  - `test_d` – tests linked to courses.
  - `question_d` – questions per test with ordering.
  - `answer_d` – answer options with scores and ordering.
  - `test_attempt_b` – student attempts and scored results.
- A shared trigger function `update_updated_at_column()` and per-table triggers keep `updated_at` in sync on updates.

The Go `adminPanel` service mirrors the `knowledge_base` structures (Go structs `Category`, `Course`, `Lesson`, etc.) and uses `github.com/jackc/pgx/v5/pgxpool` for connection pooling.

### AdminPanel API structure

- Entrypoint: `adminPanel/main.go`.
- Configuration and DB:
  - `Config` struct reads `DATABASE_URL`.
  - `initDB()` configures and pings the pgx pool (`dbPool`).
- Fiber application:
  - Uses `fiber.New` with custom JSON encoder/decoder.
  - Global middleware: CORS (`github.com/gofiber/fiber/v2/middleware/cors`) and JSON content-type enforcement.
- Routing (all JSON):
  - `/api/v1/health` – health check with DB status.
  - `/api/v1/categories` – CRUD for categories and `GET /:category_id/courses` for courses in a category.
  - `/api/v1/courses` – paginated/filterable course list and CRUD for individual courses.
  - `/api/v1/courses/:course_id/lessons` – CRUD for lessons under a course.
- Validation helpers:
  - Request-level validation of UUIDs, title lengths, allowed `level` and `visibility` values (`isValidLevel`, `isValidVisibility`, `isValidUUID`).
- Swagger integration:
  - `/swagger.json` serves `adminPanel/docs/swagger.json` with runtime-updated `basePath` and `host`.
  - `/swagger/*` is powered by `github.com/gofiber/swagger` and points at `/swagger.json`.

### Other services

- `publicSide/main.go`:
  - Minimal Go Fiber v3 app returning `"Hello, World!"` on `GET /`.
  - Useful as a starting point for building the public UI/API layer.
- `personal-account/main.py`:
  - Minimal FastAPI app returning a JSON `{"message": "Hello World"}` on `GET /`.
  - Intended to evolve to cover the `personal_account` schema.
- `testing/src/main/java/com/example/Main.java`:
  - Javalin app exposing a root endpoint and `/api/hello` JSON response.

## Testing status

- There are currently no unit test suites or framework-specific test runners configured in the repository (e.g. no Go `_test.go` files, no FastAPI test package, no JUnit tests).
- As a result, there is no established command for running a single automated test; adding tests should follow the standard patterns for the respective stacks (Go testing, pytest, JUnit) and can be wired into the existing Docker/Maven/Go tooling as needed.
