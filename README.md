# Go ORM

[![Build Status](https://travis-ci.org/posener/orm.svg?branch=master)](https://travis-ci.org/posener/orm)
[![codecov](https://codecov.io/gh/posener/orm/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/orm)
[![GoDoc](https://godoc.org/github.com/posener/orm?status.svg)](http://godoc.org/github.com/posener/orm)
[![Go Report Card](https://goreportcard.com/badge/github.com/posener/orm)](https://goreportcard.com/report/github.com/posener/orm)

An attempt to build a *typed* ORM package in Go

This is WORK IN PROGRESS

Very very limited implementations, and plenty of TODOs :-)

## Install:

```bash
go get -u github.com/posener/orm
```

## Usage

```bash
orm -h
```

The concept of this tool is to generate typed functions for a given Go struct.

## Example:

Running the orm command on the `Person` struct in the `example` package, will create a `personorm` package, with
ORM functions for the given struct.

By doing so, and having a database engine, `db`, one could do database operations with
ORM semantics.

Notice that all operations are typed, `Age` is `int`, `Name` is `string`, the `example.Person`
is used in the arguments and in the return values.

```go
import porm "package/personorm"

func main() {
    // Create table:
    err = porm.Create().Exec(db)

    // Insert rows:
    err = porm.Insert().Name("John").Age(1).Exec(db)

    // Or with a struct:
    porm.Insert().Person(&example.Person{Name: "zvika", Age: 3}).Exec(db)

    // Select rows from the DB:
    q := porm.Query().
        Select(porm.Select().Name()).
        Where(porm.WhereName(porm.OpNe, "John"))
    ps, err := q.Exec(db) // returns []example.Person, typed return value.
    println(ps[0].Name) // Output: "John"
}
```