package example

//go:generate ../orm -type A

type A struct {
	ID       int64 `sql:"primary key;auto increment"`
	Name     string
	Age      int
	CPointer *C
}

//go:generate ../orm -type B

type B struct {
	ID        int64 `sql:"primary key;auto increment"`
	Name      string
	Hobbies   string
	CsPointer []*C
}

//go:generate ../orm -type C

type C struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
	Year int
	BID  int64 `sql:"foreign key:./.B;null"`
}

//go:generate ../orm -type A2

type A2 struct {
	ID int64 `sql:"primary key;auto increment"`
	B  B2
}

//go:generate ../orm -type B2

type B2 struct {
	ID   int64 `sql:"primary key;auto increment"`
	Name string
}
