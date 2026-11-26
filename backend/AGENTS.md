# AGENTS.md — backend (bank-service)

Go 1.16 module `bank-service`. No generics or newer stdlib APIs — keep code 1.16-compatible. Development happens inside the `app` Docker container, not on the host: prefix Go commands with `docker-compose exec app`.

## Setup (order matters)

1. `cp .env.example .env` and fill the empty secrets: `BANK_SERVICE_POSTGRESQL_PASSWORD`, `JWT_SECRET_KEY`, `PGADMIN_DEFAULT_PASSWORD`. `BANK_SERVICE_USER1_PASSWORD` must be exactly 8 characters; `BANK_SERVICE_USER1_EMAIL` must be a valid email (migrations seed this user).
2. `docker network create bank-net` (external network required by compose).
3. `docker-compose build && docker-compose up`
4. Run migrations (below). Health check: `curl http://localhost:3000/ping`.

The dev container runs CompileDaemon: saving a `.go` file auto-rebuilds and restarts the server.

## Migrations

Custom runner (`robinjoseph08/go-pg-migrations`), not golang-migrate. Migration files live in `migrations/internal/` (schema + seed data) and `migrations/external/`.

```bash
docker-compose exec app go run migrations/internal/*.go migrate    # apply
docker-compose exec app go run migrations/internal/*.go rollback   # roll back last batch
docker-compose exec app go run migrations/internal/*.go create <name>
```

Migrations use **go-pg v9** while application code uses **GORM** — both ORMs coexist intentionally. New repositories use GORM.

## Tests

Tests hit a real Postgres (the compose `database` service must be up). `APP_ENV=testing` switches to the test database (`BANK_SERVICE_POSTGRESQL_NAME_TEST`, default `postgres`). `-p 1` is required: packages share one test DB and must not run in parallel.

```bash
docker-compose exec -e APP_ENV=testing app go test -cover -p 1 ./...
# single test:
docker-compose exec -e APP_ENV=testing app go test -p 1 -run TestName ./src/environments/client/modules/user/
# force re-run (Go caches results):
docker-compose exec app go clean -testcache
```

The README's `gotestsum` variant is stale — gotestsum is not installed in the image (commented out in `Dockerfile.multi-stage`).

`./audit.sh` runs a Nancy dependency-vulnerability scan (needs Docker + jq on host).

## Architecture

- Flow: `main.go` → `config.SetupCommonDependencies()` (logger, GORM DB) → `src.SetupHandler()` (middlewares + routes).
- Routes: `/v1/public/client/*` (no auth) and `/v1/client/*` (JWT-protected), wired in `src/main.go`.
- Feature modules: `src/environments/client/modules/{user,movement,contact,jwt}/`, with per-layer files `<name>.router.go` → `.controller.go` → `.service.go` → `.repository.gorm.go` (tests: `.repository.gorm_test.go`). Not every module has every layer (`jwt` is service-only). Shared entities/interfaces/routers in `src/environments/client/resources/`.
- Cross-cutting code in `src/libs/` (env, database, http, errors, middleware, i18n, ...). Read env vars through `src/libs/env`, not `os.Getenv` directly.
- Responses are localized: message catalogs in `src/libs/i18n/en.json` and `es.json` — add new keys to **both** files.
- Seed/fake users for manual testing: `fixture/fake_clients_data.json` (loaded via migrations).