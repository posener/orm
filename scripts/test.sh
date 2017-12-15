#!/usr/bin/env bash

# run project tests
# this script should be run from the project main directory

set -e

if [ -z "${TRAVIS}" ]; then
	go test ./... "$@"
else
	echo "" > coverage.txt
	for d in $(go list ./... | grep -v vendor); do
		go test -v -race -coverprofile=profile.out -covermode=atomic $d
		if [ -f profile.out ]; then
			cat profile.out >> coverage.txt
			rm profile.out
		fi
	done
fi
