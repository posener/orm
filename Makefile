all: test

.PHONY: b0x orm gen-tests test

b0x:
	go generate ./gen/...

orm: b0x
	go build ./cmd/orm

gen-tests: orm
	go generate ./tests/... ./examples/...

test: gen-tests
	./scripts/test.sh

clean:
	rm -rf tests/*_orm.go
