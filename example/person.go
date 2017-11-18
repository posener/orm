package example

//go:generate orm -name Person

type Person struct {
	Name       string
	Age        int
	unexported bool
}
