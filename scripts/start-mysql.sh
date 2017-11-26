#!/usr/bin/env bash
set -e

echo ">> Starting mysql..."
docker run --name mysql -e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} -d -p ${MYSQL_PORT}:3306 mysql:8.0.3 > /dev/null

printf ">> Waiting for mysql "
until docker exec mysql mysql -h127.0.0.1 -uroot -psecret -e "SHOW DATABASES" > /dev/null
do
    printf .
    sleep 2
done
echo ""

echo ">> Creating database"
docker exec mysql mysql -hlocalhost -uroot -psecret -e "CREATE DATABASE IF NOT EXISTS ${DB_DATABASE}" > /dev/null

echo ">> Done"
