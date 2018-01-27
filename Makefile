all: test

.PHONY: orm gen-tests test

orm:
	go build ./cmd/orm

gen-tests: orm
	go generate ./tests/... ./examples/...

test: gen-tests
	./scripts/test.sh

clean:
	find . -name "*_orm.go" -delete
