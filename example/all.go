package example

//go:generate orm -name All

// All is to test generation of variant fields and types
type All struct {
	Int        int    `sql:primary_key`
	Text       string `sql:"type:VARCHAR(100)"`
	Bool       bool
	unexported int
}
