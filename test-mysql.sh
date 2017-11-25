#!/usr/bin/env bash

echo ">> Starting mysql container"
docker run --rm --name mysql -e MYSQL_ROOT_PASSWORD=secret -d -p 3306:3306 mysql:8.0.3

set -e

echo ">> Resetting database"
docker exec mysql mysql -hlocalhost -uroot -psecret -e "DROP DATABASE IF EXISTS test; CREATE DATABASE test"

echo ">> Running tests..."
export MYSQL_ADDR="root:secret@/test"
./test.sh "$@"
