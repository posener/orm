# Go ORM

[![Build Status](https://travis-ci.org/posener/orm.svg?branch=master)](https://travis-ci.org/posener/orm)
[![GoDoc](https://godoc.org/github.com/posener/orm?status.svg)](http://godoc.org/github.com/posener/orm)
[![Go Report Card](https://goreportcard.com/badge/github.com/posener/orm)](https://goreportcard.com/report/github.com/posener/orm)

An attempt to write a **typed** ORM for Go.

> This is a PROOF OF CONCEPT

> Very very (very) limited implementations, and plenty of TODOs :-)

Check out the [Wiki](https://github.com/posener/orm/wiki) for documentation.

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
	"log"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
	porm "package/personorm"
)

func main() {
    db, err := porm.Open(dialect, source)
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
    
    // Get first matching row or "not found" error
    p, err := db.Select().First()
    
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

Simple use case is to run `orm -pkg mypackage -name MyStruct

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

See [benchmarking results](./bench).
