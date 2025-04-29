# bank-service
Basic bank service with users, movements and frequent contacts

## Development containers
In the project there is a docker-compose file used to lift the entire work environment without installing frames or tools used by the project

- Install Docker `sudo apt-get install docker-compose` (ubuntu, [for Windows use this](https://docs.docker.com/desktop/install/windows-install))

- Create the `.env` file based on `.env.example`. Use random values for secret envs

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

- Run database migrations
This step is for creating all the database tables with their fields and relations

```bash
docker-compose exec app go run migrations/internal/*.go migrate
```

You should see something like this:

```
Running batch 1 with 9 migration(s)...
Finished running "20220727011546_create-user-table"
Finished running "20220729143904_create-movement-table"
Finished running "20220730143904_trigger_set_balance_movement"
Finished running "20220730153904_trigger_create_bonus_movement"
Finished running "20220829231652_add_user1"
Finished running "20221207011102_add_movements_user1"
Finished running "20221216013623_add_movements_user1"
Finished running "20230713013541_add_fake_users"
Finished running "20230713130922_create_contact_table"
```

> If you get any error or you don't see the previous result, you can
> use:
> 
> ```bash
> 
> docker-compose  exec  app  go  run  migrations/internal/main.go migrate
> 
> ```
> 
> Or you can use these two commands:
> 
> ```bash
> 
> docker-compose  exec  app  bash
> 
> ```
> 
> ```bash
> 
> go  run  migrations/internal/*.go migrate
> 
> ```

### Congratulations! You have successfully done the configuration of this API.

## Postman docs

- [Postman documentation](https://www.postman.com/jchiquinvdev/workspace/lab3)
You should use the `Bank-service` collection, you will probably need to fork that collection into your own workspace in order to use the requests.

## Ports

|Tool            |Host                           |
|----------------|-------------------------------|
|Backend		 |`http://localhost:3000`        |
|Database        |`http://localhost:5432`		 |
|Pgadmin         |`http://localhost:80`			 |

## Fake users

You can find fake users to test this API in [the next JSON file](https://github.com/JChiquin/basic-bank-app/blob/main/backend/fixture/fake_clients_data.json)

## Migrations

- Run all internal migrations

Development
```bash
docker-compose exec app go run migrations/internal/main.go migrate
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