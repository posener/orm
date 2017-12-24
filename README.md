# Go ORM

[![Build Status](https://travis-ci.org/posener/orm.svg?branch=master)](https://travis-ci.org/posener/orm)
[![GoDoc](https://godoc.org/github.com/posener/orm?status.svg)](http://godoc.org/github.com/posener/orm)
[![Go Report Card](https://goreportcard.com/badge/github.com/posener/orm)](https://goreportcard.com/report/github.com/posener/orm)

An attempt to write a **typed** ORM for Go.

Check out the [Wiki](https://github.com/posener/orm/wiki) for documentation.

This repository gives a command line tool, called `orm`, for generating
ORM code for a given struct. The generated code is typed and has no `interface{}`s arguments
and returns values as in other ORM Go libraries.

In addition to compile time checks, editor completion and minimizing coding confusion, the generated
code results in accelerated performance - no reflections are needed in run time, and database data
decoding is faster when you know the target type.

## Example:

Check out the [examples](./examples).

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
	// OpenPersonORM was generated with the `orm` command line.
	// It retunes an object that interacts with the database.
	// This was generated for a struct:
	// type Person struct {
	// 	Name       string
	// 	Age        int
	// }
	db, err := OpenPersonORM("sqlite3", ":memory:")
	defer db.Close()
	
	// Set a logger to log SQL commands
	db.Logger(log.Printf)
	
	// Create a table:
	// The AutoMigrate modifier will cause the create table function to try and
	// migrate an existing table. Current implementation supports adding columns
	// and foreign keys.
	err = db.Create().AutoMigrate().Exec()
	
	// Insert a row with arguments:
	// The Insert() function returns a builder, which has functions according
	// to the struct's fields. The nice thing here: SetName's argument is of type
	// string, and SetAge's argument is of type int.
	// No room for run time errors!
	// The returned values are the created object, of type *Person, and an error,
	// in case that there was an error in the INSERT operation. Everything is typed!
	john, err = db.Insert().SetName("John").SetAge(1).Exec()
	println(john.Name) // Output: John

	// Insert a row with a struct:
	// InsertPerson's argument is of type *Person, no room for run tme errors!
	// Again, Exec() returns a *Person and error objects.
	doug, err = db.Insert().InsertPerson(&tests.Person{Name: "Doug", Age: 3}).Exec()
	println(doug.Name, doug.Age) // Output: Doug 3

	// Select rows from the table:
	// The Where() modifier adds a WHERE statement to the SQL query, it's input is a 
	// generated function with generated modifiers. Name() accepts a constant that 
	// defines the operation, and the second argument is of type string, so only valid
	// comparisons can be made, and they are checked in compile time.
	// Query() returns the slice []Person, and an error object. Everything is typed!
	persons, err := db.Select().
		Where(db.Where().Name(orm.OpNe, "John")).
		Query() // returns []tests.Person, typed return value.
	println(persons[0].Age) // Output: 1
	
	// Get first matching row or "not found" error
	// Again, the first argument is of type *Person!
	person, err := db.Select().First()
	
	// Delete a row:
	// Here, we are using the Age() function of the Where() modifier, and it's second
	// argument is of type int, so only valid comparisons can be made.
	_, err = db.Delete().Where(db.Where().Age(orm.OpGt, 12)).Exec()
	
	// I think you get the vibe by now, needless to say it again. Everything is ...
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

> [GORM](http://jinzhu.me/gorm/): is a mature, widely used, heavily debugged, favorable, 
  and very recommended to use. `orm` advantages over GORM are: fully typed APIs, clear and explicit APIs, speed,
  support for `Context`, and support for custom logger.

Also, [benchmarking results](./bench) are available.
