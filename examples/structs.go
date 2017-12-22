package examples

//go:generate ../orm -type Simple,One,OtherOne,OtherMany

type Simple struct {
	ID     int64 `sql:"primary key;auto increment"`
	Field1 string
	Field2 int
	Field3 bool
}

type One struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	// OtherOne is one-to-one relationship with a type called "OtherOne"
	OtherOne *OtherOne
	// OtherMany is one-to-many relationship with a type called "OtherMany"
	// In order for this relationship to exists, OtherMany must have a
	// field that references a "One" type
	OtherMany []OtherMany
}

type OtherOne struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
}

type OtherMany struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	// MyOneIs is a filed that must exists for allowing a one to many relationship
	// between a One type and OtherMany type.
	MyOne *One
}
