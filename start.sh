#!/bin/sh

set -e

echo "Run DB migrations"
source config.env
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Start Application"
# Dockerfile: When CMD is used with ENTRYPOINT, it acts as additional parameters to ENTRYPOINT.
# Making it run as: /app/start.sh /app/main
# exec "$@" executes the parameter /app/main
exec "$@"