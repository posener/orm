#! /usr/bin/env bash

set -e
cd "$(dirname "$0")"
go test -bench . > benchmark.txt
go run ./plot.go

