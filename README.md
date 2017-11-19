# Go ORM

[![Build Status](https://travis-ci.org/posener/orm.svg?branch=master)](https://travis-ci.org/posener/orm)
[![codecov](https://codecov.io/gh/posener/orm/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/orm)
[![GoDoc](https://godoc.org/github.com/posener/orm?status.svg)](http://godoc.org/github.com/posener/orm)
[![Go Report Card](https://goreportcard.com/badge/github.com/posener/orm)](https://goreportcard.com/report/github.com/posener/orm)

An proof of concept for a Go **typed** ORM package.

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
	porm "package/personorm"
)

func main() {
    db, err := sql.Open(...)
    defer db.Close()

    // Create table:
    err = porm.Create().Exec(db)

    // Insert rows:
    err = porm.Insert().SetName("John").SetAge(1).Exec(db)

    // Or with a struct:
    porm.InsertPerson(&example.Person{Name: "Doug", Age: 3}).Exec(db)

    // Select rows from the DB:
    ps, err := porm.Query().
    	SelectAge().
        Where(porm.WhereName(porm.OpNe, "John")).
        Exec(db) // returns []example.Person, typed return value.

    println(ps[0].Age) // Output: 1
    
    // Delete
    porm.Delete().Where(porm.WhereName(porm.Eq, "John")).Exec(db)
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
