package example

//go:generate ../orm -type A,B,C,A2,B2,C2,D2,A3,B3

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
	BID  int64 `sql:"foreign key:./.B;null"`
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
