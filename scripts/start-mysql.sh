#!/usr/bin/env bash
set -e

DB_PASSWORD=${DB_PASSWORD:-secret}
MYSQL_PORT=${MYSQL_PORT:-3306}
DB_DATABASE=${DB_DATABASE:-test}

echo ">> Starting mysql..."
docker run --name mysql -e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} -d -p ${MYSQL_PORT}:3306 mysql:5.6.38 > /dev/null

printf ">> Waiting for mysql "
until docker exec mysql mysql -h127.0.0.1 -uroot -psecret -e "SHOW DATABASES" > /dev/null
do
    printf .
    sleep 2
done
echo ""

echo ">> Creating database"
docker exec mysql mysql -hlocalhost -uroot -p${DB_PASSWORD} -e "CREATE DATABASE IF NOT EXISTS ${DB_DATABASE}" > /dev/null

echo ">> Done

To run mysql tests, export the following environment variable:

export MYSQL_ADDR=root:secret@(127.0.0.1:${MYSQL_PORT})/${DB_DATABASE}

"

