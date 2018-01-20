package examples

//go:generate ../../orm -type One,OtherOne

// One is for testing ont-to-one relation
type One struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	// OtherOne is one-to-one relationship with a type called "OtherOne"
	OtherOne *OtherOne
}

// OtherOne is also for testing one-to-one relation
type OtherOne struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
}
