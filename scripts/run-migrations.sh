#!/usr/bin/env bash

. ./scripts/env.sh

docker run -v --rm \
-v ./migrations:/migrations \
--network ${NAME}_open_fermentations \
migrate/migrate \
-path=/migrations/ \
-database ${DB_CONNECTION_STRING} \
up