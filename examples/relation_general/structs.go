package relationone2many

//go:generate ../../orm -type One,Two

// One is a struct for example
type One struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	// A relation to Two, since Two doesn't have a field which is a pointer to One,
	// This is considered to be a general relation, which will be stored in a table
	// called `rel_one_two` if created.
	Two []Two
	// A relation to Two with a "relation name" modifier, the relation table will be
	// called "mutual" since this is the relation name.
	// Since Two also has a field with the same relation name, this relation is mutual,
	// and adding a relation with One's orm will be shown in Two's queries and vice versa.
	Mutual []Two `sql:"relation name:mutual"`
}

// Two is a struct for example
type Two struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	// A relation to One, since One doesn't have a field which is a pointer to Two,
	// This is considered to be a general relation, which will be stored in a table
	// called `rel_two_one` if created.
	One []One
	// A relation to One with a "relation name" modifier, the relation table will be
	// called "mutual" since this is the relation name.
	// Since One also has a field with the same relation name, this relation is mutual,
	// and adding a relation with Two's orm will be shown in One's queries and vice versa.
	Mutual []One `sql:"relation name:mutual"`
}
