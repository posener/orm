package example

//go:generate orm -name Person -dialect sqlite3

// Person is en example struct for the ORM package
type Person struct {
	Name       string
	Age        int
	unexported bool
}
