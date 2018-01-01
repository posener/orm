package tests

import "time"

//go:generate ../orm -type Person,Employee,All
//go:generate ../orm -type Person -out ./external

// Person for testing
type Person struct {
	Name       string
	Age        int
	unexported bool
}

// Employee is a person who works
// This is a test case for struct embedding
type Employee struct {
	Person
	Salary int
}

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

// Ignore contains ignored fields.
type Ignore struct {
	ID   int64
	Data map[string]interface{} `sql:"-"`
}
