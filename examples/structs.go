package examples

//go:generate ../orm -type Simple

type Simple struct {
	ID     int64 `sql:"primary key;auto increment"`
	Field1 string
	Field2 int
	Field3 bool
}
