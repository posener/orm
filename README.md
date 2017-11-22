# Go ORM

[![Build Status](https://travis-ci.org/posener/orm.svg?branch=master)](https://travis-ci.org/posener/orm)
[![GoDoc](https://godoc.org/github.com/posener/orm?status.svg)](http://godoc.org/github.com/posener/orm)
[![Go Report Card](https://goreportcard.com/badge/github.com/posener/orm)](https://goreportcard.com/report/github.com/posener/orm)

An attermpt to write a **typed** ORM for Go.

> This is a PROOF OF CONCEPT

> Very very (very) limited implementations, and plenty of TODOs :-)

This repository gives a command line tool, called `orm`, for generating
ORM code for a given struct. The generated code is typed and has no `interface{}`s arguments
and return values as in other ORM Go libraries.

## Example:

Running the orm command on the `Person` struct in the `example` package, will create a `personorm` package, with
ORM functions for the given struct.

By doing so, and having a database engine, `db`, one could do database operations with
ORM semantics.

Notice that all operations are typed, `Age` is `int`, `Name` is `string`, the `example.Person`
is used in the arguments and in the return values.

```go
import (
	"database/sql"
	"log"
	
	porm "package/personorm"
)

func main() {
    db, err := sql.Open(...)
    defer db.Close()
    
    orm = porm.New(db)
    
    // Set a logger to log SQL commands
    orm.Logger(log.Printf)

    // Create table:
    _, err = orm.Create().Exec()

    // Insert row with arguments:
    _, err = orm.Insert().SetName("John").SetAge(1).Exec()

    // Insert row with a struct:
    _, err = orm.InsertPerson(&example.Person{Name: "Doug", Age: 3}).Exec()

    // Select rows from the table:
    ps, err := orm.Select().
    	SelectAge().
        Where(porm.WhereName(porm.OpNe, "John")).
        Query() // returns []example.Person, typed return value.

    println(ps[0].Age) // Output: 1
    
    // Delete row
    _, err = orm.Delete().Where(porm.WhereName(porm.Eq, "John")).Exec()
}
```

## Command Line Tool

### Install:

```bash
go get -u github.com/posener/orm
```

### Usage

Run `orm -h` to get detailed usage information.

Simple use case is to run `orm -pkg mypackage -name MyStruct`.

#### go generate

By adding the comment aside to the type deceleration, as shown below, one could run `go generate ./...`
to generate the ORM files for `MyType`.

```go
//go:generate orm -name MyType

type MyType struct {
	...
}
```

## Benchmark

Initial benchmark tests are available in [here](./example/bench_test.go).

#### Compared packages:

- [x] posener/orm (this package)
- [x] jinzhu/gorm
- [x] Direct SQL commands

#### Operations:

- [x] INSERT

### Results:

```bash
$ go test -bench . -benchtime 10s
goos: linux
goarch: amd64
pkg: github.com/posener/orm
BenchmarkORM-4    	 1000000	     21229 ns/op
BenchmarkGORM-4   	  200000	     63628 ns/op
BenchmarkRaw-4    	 1000000	     18283 ns/op
PASS
ok  	github.com/posener/orm	53.333s
```

# Documentation

## Model definition

In the model, it is possible to annotate a field with Go tags.
`orm` will consider the `sql` keyword in the tag of any field.

### Supported Tags

The supported tags are shown in [ParseTags](./gen/tags.go) function.

- `type`: Type of SQL column.
- `primary_key`/`primary key`: Is the field a primary key. (can be given to multiple fields in the same struct)
- `not null`/`not_null`: Field is not allowed to be null.
- `auto_increment`/`auto increment`: Apply auto increment for this column.
- `unique`: Each row in this column must be unique.
- `default`: Give a default value for this field.
