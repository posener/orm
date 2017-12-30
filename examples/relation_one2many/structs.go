package relation_one2many

//go:generate ../../orm -type One,OtherMany

type One struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	// OtherMany is one-to-many relationship with a type called "OtherMany"
	// In order for this relationship to exists, OtherMany must have a
	// field that references a "One" type
	OtherMany []OtherMany
}

type OtherMany struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	// MyOneIs is a filed that must exists for allowing a one to many relationship
	// between a One type and OtherMany type.
	MyOne *One
}
