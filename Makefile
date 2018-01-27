all: test

.PHONY: orm gen-tests test bench profile

orm:
	go build ./cmd/orm

gen-tests: orm
	go generate ./tests/... ./examples/...

test: gen-tests
	./scripts/test.sh

clean:
	find . -name "*_orm.go" -delete

bench:
	cd bench && go test -bench . > benchmark.txt
	cd bench && go run ./plot.go

profile:
	cd bench && go test -bench BenchmarkORM -cpuprofile cpu.out -memprofile mem.out
	go tool pprof -web bench/cpu.out
	go tool pprof -web bench/mem.out
