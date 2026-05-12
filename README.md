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

## MQTT

MVP setup could be as simple as having a pwfile that is already set up with 2 users in it. One for the platform to be able to interact with the server and then one for the devices.

Ideally every user on the system should have their own devices user but that would require us to generate the user and then restart the mqtt server for that to take effect.
This is not ideal.

It seems like there are some plugins that could help with this like dynamic authentication or something.

### Setup

#### Create pwfile

Add or change the `MQTT_PASSWORD` in [.env](.env)

```bash
cp .mqtt/config/mosquitto.conf.example .mqtt/config/mosquitto.conf
docker compose up mqtt -d
docker compose exec mqtt sh # opens the terminal in the docker container
mosquitto_passwd -c /mosquitto/config/pwfile platform
# use the password that you set in the .env
chmod 700 /mosquitto/config/pwfile # sets read write and execute for the owner
chown mosquitto:mosquitto /mosquitto/config/pwfile # sets the mosquitto user and group as the owner of the file
exit # exits the terminal in the docker container
```

#### Enable pwfile

In [pwfile](./.mqtt/config/mosquitto.conf)

- Disable anonymous authentication
- Uncomment line to point to the password file

```bash
...
allow_anonymous false
password_file /mosquitto/config/pwfile
...
```

After changing the config restart the container for the changes to take effect

```bash
docker compose restart mqtt
```