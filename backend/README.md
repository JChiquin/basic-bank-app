# bank-service
Basic bank service with users, movements and frequent contacts

## Development containers
In the project there is a docker-compose file used to lift the entire work environment without installing frames or tools used by the project. It only requires [Docker](https://docs.docker.com/get-docker/) with the `docker compose` plugin (v2).

- Create the `.env` file based on `.env.example`. Use random values for secret envs

- Build docker orchestra

```bash
docker compose build
```

- Run docker orchestra

```bash
docker compose up
```

The database container is waited until it is **healthy**, then the app container:
1. Runs all pending database migrations automatically (idempotent, safe on every start).
2. Starts the API with hot-reload: editing any `.go` file rebuilds and restarts the service.

- Check service up

```bash
curl http://localhost:3000/ping
```

No external networks to create, no manual migration step.

### pgAdmin (optional)
pgAdmin is not started by default to keep the environment light. Start it only when you want a database UI:

```bash
docker compose --profile tools up pgadmin
```

Then open http://localhost:5050 (credentials are `PGADMIN_DEFAULT_EMAIL` / `PGADMIN_DEFAULT_PASSWORD` from your `.env`) and register the server with host `bank-service-postgres`, port `5432` and your database credentials (`BANK_SERVICE_POSTGRESQL_USERNAME` and `BANK_SERVICE_POSTGRESQL_PASSWORD`).

If you prefer a local client (DBeaver, psql, etc.), uncomment the `5432` port mapping in `docker-compose.yml` and connect to `localhost:5432`.

## Postman docs

- [Postman documentation](https://www.postman.com/jchiquindev/lab-ii)
You should use the `Bank-service` collection, you will probably need to fork that collection into your own workspace in order to use the requests.

## Ports

|Tool            |Host                           |
|----------------|-------------------------------|
|Backend		 |`http://localhost:3000`        |
|pgAdmin (optional)|`http://localhost:5050`		 |

The database is not exposed to the host by default; uncomment its port mapping in `docker-compose.yml` if you need it.

## Fake users
You can find fake users to test this API in [the next JSON file](https://github.com/JChiquin/basic-bank-app/blob/main/backend/fixture/fake_clients_data.json)

## Migrations

Migrations run automatically when the app container starts, before the server boots. Two manual commands remain available inside the running container:

```bash
# Roll back the previous batch of migrations
docker compose exec app go run migrations/internal/*.go rollback

# Create a new migration file
docker compose exec app go run migrations/internal/*.go create <migration_name>
```

External migrations live in `migrations/external/` and follow the same pattern.

## Unit tests

Before the first run, create the tables in the test database (they use `BANK_SERVICE_POSTGRESQL_NAME_TEST`):

```bash
docker compose exec -e APP_ENV=testing app go run migrations/internal/*.go migrate
```

- Run all test (`-p 1` is required: packages share one test database and must not run in parallel)

```bash
docker compose exec -e APP_ENV=testing app go test -cover -p 1 ./...
```

- Run a single test

```bash
docker compose exec -e APP_ENV=testing app go test -p 1 -run TestName ./src/environments/client/modules/user/
```

- Clean cache (Useful for forcing tests to run again)

```bash
docker compose exec app go clean -testcache
```
