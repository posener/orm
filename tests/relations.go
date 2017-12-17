package tests

//go:generate ../orm -type A,B,C,A2,B2,C2,D2,A3,B3,A4,B4,A5,A6,B6,A7,B7

// A, B, C test simple one-to-one (A->C) and one-to-many(B->C) relationships
type A struct {
	ID       int64 `sql:"primary key;auto increment"`
	Name     string
	Age      int
	CPointer *C
}

type B struct {
	ID        int64 `sql:"primary key;auto increment"`
	Name      string
	Hobbies   string
	CsPointer []*C
}

type C struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	Year int
	B    *B
}

// A2, B2, C2, D2 test nested relationship (A2->B2->C2,D2),
// and non-pointer reference(A2->B2)
type A2 struct {
	ID int64 `sql:"primary key;auto increment"`
	B  B2
}

type B2 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	C    *C2
	D    *D2
}

type C2 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
}

type D2 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
}

// A3, B3 test bi-directional one-to-many relation
type A3 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	B    []*B3
}

type B3 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	A    *A3 `sql:"foreign key:./.A3;null"`
}

// A4, B4 test multiple fields with the same reference type
type A4 struct {
	ID     int64 `sql:"primary key;auto increment"`
	Name   string
	B1, B2 *B4
}

type B4 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
}

// A5 test self referencing
type A5 struct {
	ID          int64 `sql:"primary key;auto increment"`
	Name        string
	Left, Right *A5
}

// A6,B6 test referencing unique keys
type A6 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	B    *B6
}

type B6 struct {
	SureName  string `sql:"primary key"`
	FirstName string `sql:"primary key"`
}

// A7,B7 test functionality of 'referencing field'
// A7 has one-to-many relationship to B7, and B7 has several A7 reverse references
type A7 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	B    []B7 `sql:"referencing field:A1"`
}

type B7 struct {
	ID     int64 `sql:"primary key;auto increment"`
	Name   string
	A1, A2 *A7
}
