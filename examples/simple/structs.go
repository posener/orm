package simple

//go:generate ../../orm -type Simple

// Simple is a simple struct for example
type Simple struct {
	ID     int64 `sql:"primary key;auto increment"`
	Field1 string
	Field2 int
	Field3 bool
}
