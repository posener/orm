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

Running the orm command on the `Person` struct in the `example` package with `sqlite3` dialect, 
will create a `personsqlite3` package, with ORM functions for the given struct.

By doing so, and having a database engine, `db`, one could do database operations with
ORM semantics.

Notice that all operations are typed, `Age` is `int`, `Name` is `string`, the `example.Person`
is used in the arguments and in the return values.

```go
import (
	"database/sql"
	"log"
	
	"github.com/posener/orm"
	porm "package/personorm"
)

func main() {
    db, err := porm.Open(source)
    defer db.Close()
    
    
    // Set a logger to log SQL commands
    db.Logger(log.Printf)

    // Create table:
    _, err = db.Create().Exec()

    // Insert row with arguments:
    _, err = db.Insert().SetName("John").SetAge(1).Exec()

    // Insert row with a struct:
    _, err = db.InsertPerson(&example.Person{Name: "Doug", Age: 3}).Exec()

    // Select rows from the table:
    ps, err := db.Select().
    	SelectAge().
        Where(porm.WhereName(orm.OpNe, "John")).
        Query() // returns []example.Person, typed return value.

    println(ps[0].Age) // Output: 1
    
    // Delete row
    _, err = db.Delete().Where(porm.WhereName(orm.Eq, "John")).Exec()
}
```

## Command Line Tool

### Install:

```bash
go get -u github.com/posener/orm/cmd/orm
```

### Usage

Run `orm -h` to get detailed usage information.

Simple use case is to run `orm -pkg mypackage -name MyStruct -dialect sqlite3`.

#### go generate

By adding the comment aside to the type deceleration, as shown below, one could run `go generate ./...`
to generate the ORM files for `MyType`.

```go
//go:generate orm -name MyType -dialect sqlite3

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
- [X] SELECT
- [X] SELECT with big structs

### Results:

```bash
$ go test -bench . -benchtime 10s
goos: linux
goarch: amd64
pkg: github.com/posener/orm
BenchmarkORMInsert-4              	 1000000	     21634 ns/op
BenchmarkGORMInsert-4             	  200000	     67917 ns/op
BenchmarkRawInsert-4              	 1000000	     18546 ns/op
BenchmarkORMQuery-4               	    5000	   2954333 ns/op
BenchmarkGORMQuery-4              	    2000	   9880517 ns/op
BenchmarkRawQuery-4               	    5000	   2672430 ns/op
BenchmarkORMQueryLargeStruct-4    	    2000	   6741345 ns/op
BenchmarkGORMQueryLargeStruct-4   	    1000	  20876834 ns/op
BenchmarkRawQueryLargeStruct-4    	    2000	   7429866 ns/op
PASS
ok  	github.com/posener/orm	167.655s

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
