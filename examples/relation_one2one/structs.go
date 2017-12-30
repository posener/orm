package examples

//go:generate ../../orm -type One,OtherOne

type One struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	// OtherOne is one-to-one relationship with a type called "OtherOne"
	OtherOne *OtherOne
}

type OtherOne struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
}
