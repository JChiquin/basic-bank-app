# AGENTS.md — backend (bank-service)

Go 1.23 module `bank-service` (dockerized toolchain: `golang:1.23-alpine`, pure-Go builds with `CGO_ENABLED=0`). Development happens inside the `app` Docker container, not on the host: prefix Go commands with `docker compose exec app` (Compose v2 syntax; `docker-compose` v1 is not supported by this file).

## Setup

1. `cp .env.example .env` and fill the empty secrets: `BANK_SERVICE_POSTGRESQL_PASSWORD`, `JWT_SECRET_KEY`, `PGADMIN_DEFAULT_PASSWORD`. `BANK_SERVICE_USER1_PASSWORD` must be exactly 8 characters; `BANK_SERVICE_USER1_EMAIL` must be a valid email (migrations seed this user).
2. `docker compose build && docker compose up`

That's all. The app container waits for the Postgres healthcheck, **runs migrations automatically**, then starts the API with hot-reload (CompileDaemon: saving a `.go` file auto-rebuilds and restarts, ~1s with the warm module cache baked into the image). Health check: `curl http://localhost:3000/ping`. If a migration fails, the container exits with the error in `docker compose logs app` (fail-fast by design).

pgAdmin is opt-in: `docker compose --profile tools up -d pgadmin` → http://localhost:5050, connects to host `bank-service-postgres:5432`. Postgres is not published to the host unless you uncomment the port mapping in `docker-compose.yml`.

## Migrations

Custom runner (`robinjoseph08/go-pg-migrations`), not golang-migrate. Migration files live in `migrations/internal/` (schema + seed data) and `migrations/external/`. `migrate` is idempotent (state table) and runs automatically on container start. Manual commands inside the running container:

```bash
docker compose exec app go run migrations/internal/*.go rollback    # roll back last batch
docker compose exec app go run migrations/internal/*.go create <name>
```

- Migrations use **go-pg v9** while application code uses **GORM** — both ORMs coexist intentionally. New repositories use GORM.
- The fake-users migration reads `./fixture/fake_clients_data.json` at runtime with a **relative path** — it must run with cwd `/app` (the container's WORKDIR; safe in all documented commands).
- The migrations binary switches to the test database when `APP_ENV=testing`.

## Tests

Tests hit a real Postgres (the compose `database` service must be up). `APP_ENV=testing` switches to the test database (`BANK_SERVICE_POSTGRESQL_NAME_TEST`, default `postgres`). Before the first run, migrate the test DB: `docker compose exec -e APP_ENV=testing app go run migrations/internal/*.go migrate`. `-p 1` is required: packages share one test DB and must not run in parallel.

```bash
docker compose exec -e APP_ENV=testing app go test -cover -p 1 ./...
# single test:
docker compose exec -e APP_ENV=testing app go test -p 1 -run TestName ./src/environments/client/modules/user/
# force re-run (Go caches results):
docker compose exec app go clean -testcache
```

`./audit.sh` runs a Nancy dependency-vulnerability scan (needs Docker + jq on host).

## Architecture

- Flow: `main.go` → `config.SetupCommonDependencies()` (logger, GORM DB) → `src.SetupHandler()` (middlewares + routes).
- Routes: `/v1/public/client/*` (no auth) and `/v1/client/*` (JWT-protected), wired in `src/main.go`.
- Feature modules: `src/environments/client/modules/{user,movement,contact,jwt}/`, with per-layer files `<name>.router.go` → `.controller.go` → `.service.go` → `.repository.gorm.go` (tests: `.repository.gorm_test.go`). Not every module has every layer (`jwt` is service-only). Shared entities/interfaces/routers in `src/environments/client/resources/`.
- Cross-cutting code in `src/libs/` (env, database, http, errors, middleware, i18n, ...). Read env vars through `src/libs/env`, not `os.Getenv` directly. Note: the DB retry interval env var is `BANK_SERVICE_SECONDS_BETWEEN_ATTEMPTS` (the bare `SECONDS_BETWEEN_ATTEMPTS` name is ignored).
- Responses are localized: message catalogs in `src/libs/i18n/en.json` and `es.json` — add new keys to **both** files.
- Seed/fake users for manual testing: `fixture/fake_clients_data.json` (loaded via migrations).
