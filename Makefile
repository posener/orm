all: test

.PHONY: b0x orm tests test

b0x:
	go generate ./gen/...

orm: b0x
	go build ./cmd/orm

tests: orm
	go generate ./tests/...

test: tests
	./scripts/test.sh

clean:
	rm -rf tests/*_orm.go
