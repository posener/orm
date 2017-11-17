#!/usr/bin/env bash

set -e

echo "Generating Person ORM..."
go run main.go -pkg ./example -name Person

echo "Running tests..."
go test -v ./...
