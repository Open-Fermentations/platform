# Project open-fermentations

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

Generate queries for repository
```bash
make sqlc
```

## Database migrations

This project makes use of [golang migrate](https://github.com/golang-migrate/migrate).

In order to run the some scripts below you will have to install the CLI of the migration tool.

There is a service in the docker compose file that will run the migrations at start of the system to ensure that the database is up to date.
Once the migrations have run successfully the application will be able to start.

### Create

```bash
./scripts/create-migration.sh <migration-name>
```

Will create an `up` and `down` migration file in the [migrations](/migrations) folder

### Rolling back

```bash
./scripts/rollback-migration.sh
```

Rolls back the latest migration. This script can be ran until no more migrations are present