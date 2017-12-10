all: test

.PHONY: b0x orm examples test

b0x:
	go generate ./gen/...

orm: b0x
	go build ./cmd/orm

examples: orm
	go generate ./example/...

test: examples
	./test.sh

clean:
	rm -rf example/*_orm.go
