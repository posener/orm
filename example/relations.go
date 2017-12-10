package example

//go:generate ../orm -type A -type B -type C -type A2 -type B2

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

type A2 struct {
	ID int64 `sql:"primary key;auto increment"`
	B  B2
}

type B2 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
}
