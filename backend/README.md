# bank-service
Basic bank service with users and movements

## Development containers
In the project there is a docker-compose file used to lift the entire work environment without installing frames or tools used by the project

- Install Docker `sudo apt-get install docker-compose` (ubuntu, [for Windows use this](https://docs.docker.com/desktop/install/windows-install))

- Create the `.env` file based on `.env.example`.

- Create a global bridge network

```bash
docker network create bank-net
```

- Build docker orchestra

```bash
docker-compose build
```

- Run docker orchestra

```bash
docker-compose up
```

- Check service up

```bash
curl http://localhost:3000/ping
```

- Run databae migrations
```bash
docker-compose exec app go run migrations/internal/*.go migrate
```

## Ports

|Tool            |Host                           |
|----------------|-------------------------------|
|Backend		 |`http://localhost:3000`        |
|Database        |`http://localhost:5432`		 |
|Pgadmin         |`http://localhost:80`			 |

## Migrations

- Run all internal migrations

Development
```bash
docker-compose exec app go run migrations/internal/*.go migrate
```

If you need to rollback, just change the final word `migrate` for `rollback` (roll back the previous run batch of migrations)

- Create internal migrations

```bash
go run migrations/internal/*.go create <migration_name>
```

- Create external migrations

```bash
go run migrations/external/*.go create <migration_name>
```

## Unit tests

- Run all test

```bash
docker-compose exec -e APP_ENV=testing app go test -cover -p 1 ./...
```

With better format:
```bash
docker-compose exec -e APP_ENV=testing app gotestsum --format testname -- ./... -p 1 -count 1 -cover -coverprofile cover.out
```

Clean cache (Useful for forcing tests to run again)
```bash
docker-compose exec app go clean -testcache
```

## Fake users

You can find fake users to test this API in [the next JSON file](https://github.com/JChiquin/basic-bank-app/blob/main/backend/fixture/fake_clients_data.json)