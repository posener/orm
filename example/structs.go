package example

import "time"

//go:generate orm -type Person

// Person is en example struct for the ORM package
type Person struct {
	Name       string
	Age        int
	unexported bool
}

//go:generate orm -type Employee

// Employee is a person who works
// This is a test case for struct embedding
type Employee struct {
	Person
	Salary int
}

//go:generate orm -type Loaner

type Loaner struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	Age  int
	Book *Book
}

//go:generate orm -type Book

type Book struct {
	ID       int64 `sql:"primary key;auto increment"`
	Name     string
	Year     int
	AuthorID int64 `sql:"foreign key:./example.Author;null"`
}

//go:generate orm -type Author

type Author struct {
	ID      int64 `sql:"primary key;auto increment"`
	Name    string
	Hobbies string
	Books   []*Book
}

//go:generate orm -type All

// All is to test generation of variant fields and types
type All struct {
	// Annotated
	Auto   int    `sql:"primary key;autoincrement"`
	NotNil string `sql:"not null"`

	// Standard

	Int    int
	Int8   int8
	Int16  int16
	Int32  int32
	Int64  int64
	UInt   uint
	UInt8  uint8
	UInt16 uint16
	UInt32 uint32
	UInt64 uint64

	Time time.Time

	VarCharString string `sql:"type:VARCHAR(100)"`
	VarCharByte   []byte `sql:"type:VARCHAR(100)"`
	String        string
	Bytes         []byte

	Bool bool

	// Pointers

	PInt    *int
	PInt8   *int8
	PInt16  *int16
	PInt32  *int32
	PInt64  *int64
	PUInt   *uint
	PUInt8  *uint8
	PUInt16 *uint16
	PUInt32 *uint32
	PUInt64 *uint64

	PTime *time.Time

	PVarCharString *string `sql:"type:VARCHAR(100)"`
	PVarCharByte   *[]byte `sql:"type:VARCHAR(100)"`
	PString        *string
	PBytes         *[]byte

	PBool *bool

	// Special cases

	// test that unexported are not generating columns
	unexported int

	// test a case where field is a reserved name
	Select int
}
