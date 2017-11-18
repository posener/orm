package example

//go:generate orm -name All

type All struct {
	String     string `sql:"type:VARCHAR(100)"`
	Int        int
	Bool       bool
	unexported int
}
