#!/usr/bin/env bash
set -e

DB_PASSWORD=
POSTGRES_PORT=${POSTGRES_PORT:-5432}
DB_DATABASE=${DB_DATABASE:-test}
DB_USER=postgres

echo ">> Starting postgres..."
docker run --name postgres -e POSTGRES_PASSWORD=${DB_PASSWORD} -d -p ${POSTGRES_PORT}:5432 postgres:9.6.6-alpine > /dev/null

printf ">> Waiting for mysql "
until docker exec postgres psql -U ${DB_USER} -c "SELECT 1"
do
    printf .
    sleep 2
done
echo ""

echo ">> Creating database"
docker exec postgres psql -U postgres -c "CREATE DATABASE ${DB_DATABASE}" > /dev/null

echo ">> Done

To run postgres tests, export the following environment variable:

export POSTGRES_ADDR=postgres://${DB_USER}:${DB_PASSWORD}@127.0.0.1:${POSTGRES_PORT}/?sslmode=disable

"

