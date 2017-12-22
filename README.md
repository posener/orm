# Go ORM

[![Build Status](https://travis-ci.org/posener/orm.svg?branch=master)](https://travis-ci.org/posener/orm)
[![GoDoc](https://godoc.org/github.com/posener/orm?status.svg)](http://godoc.org/github.com/posener/orm)
[![Go Report Card](https://goreportcard.com/badge/github.com/posener/orm)](https://goreportcard.com/report/github.com/posener/orm)

An attempt to write a **typed** ORM for Go.

> This is a PROOF OF CONCEPT

> Very very (very) limited implementations, and plenty of TODOs :blush:

Check out the [Wiki](https://github.com/posener/orm/wiki) for documentation.

Check out the [examples](./examples).

This repository gives a command line tool, called `orm`, for generating
ORM code for a given struct. The generated code is typed and has no `interface{}`s arguments
and returns values as in other ORM Go libraries.

In addition to compile time checks, editor completion and minimizing coding confusion, the generated
code results in accelerated performance - no reflections are needed in run time, and database data decoding is faster when you know the target type.

## Example:

Running the orm command on the `Person` struct in the `tests` package with `sqlite3` dialect, 
will create a `personsqlite3` package, with ORM functions for the given struct.

By doing this, given a working database engine, `db`, you could perform database operations with
ORM semantics.

Notice that all operations are typed, `Age` is an `int`, `Name` is a `string`, the `tests.Person`
is used in the arguments and in the return values.

```go
import (
	"log"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/posener/orm"
)

func main() {
    db, err := OpenPersonORM(dialect, source)
    defer db.Close()
    
    // Set a logger to log SQL commands
    db.Logger(log.Printf)

    // Create table:
    err = db.Create().AutoMigrate().Exec()

    // Insert row with arguments:
    john, err = db.Insert().SetName("John").SetAge(1).Exec()
    println(john.Name) // Output: John

    // Insert row with a struct:
    doug, err = db.Insert().InsertPerson(&tests.Person{Name: "Doug", Age: 3}).Exec()
    println(doug.Name, doug.Age) // Output: Doug 3

    // Select rows from the table:
    persons, err := db.Select(PersonColAge).
        Where(db.Where().Name(orm.OpNe, "John")).
        Query() // returns []tests.Person, typed return value.

    println(persons[0].Age) // Output: 1
    
    // Get first matching row or "not found" error
    person, err := db.Select().First()
    
    // Delete row
    _, err = db.Delete().Where(db.Where().Name(orm.Eq, "John")).Exec()
}
```

## Command Line Tool

### Install:

```bash
go get -u github.com/posener/orm/cmd/orm
```

### Usage

Run `orm -h` to get detailed usage information.

A simple use-case is to run `orm -type MyStruct

#### go generate

By adding the comment aside to the type deceleration, as shown below, one could run `go generate ./...`
to generate the ORM files for `MyType`.

```go
//go:generate orm -type MyType

type MyType struct {
	...
}
```

### Currently Supported Dialects:

- [x] mysql
- [x] sqlite

### Compare with other Go ORMs 

Here are the strengths of `orm` in comparison to other ORM libraries for Go.

> Just to clearify: [GORM](http://jinzhu.me/gorm/) is a mature, widely used, heavily debugged, favorable, 
  and very recommended to use. `orm` is in **proof of concept** stage.

|     | `orm` | [GORM](http://jinzhu.me/gorm/) |
| --- | :---: | :----: |
| Typed API | :muscle: | :disappointed: |
| `Context` support | :muscle: | :disappointed: |
| Custom logger | :muscle: | :disappointed: |
| Speed | :rabbit: | :snail: |
| Clear and explicit API | :bowtie: | :confused: |


Also, [benchmarking results](./bench) are available.
