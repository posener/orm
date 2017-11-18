package example

//go:generate orm -name All

// All is to test generation of variant fields and types
type All struct {
	Text       string `sql:"type:VARCHAR(100)"`
	Int        int
	Bool       bool
	unexported int
}
