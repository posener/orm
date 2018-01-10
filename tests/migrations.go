package tests

//go:generate ../orm -type Migration0,Migration1,Migration2,Migration3

type Migration0 struct {
	A string
}

func (*Migration0) TableName() string {
	return "migration"
}

type Migration1 struct {
	A string
	B string
}

func (*Migration1) TableName() string {
	return "migration"
}

type Migration2 struct {
	A  string
	B  string
	D  string
	P1 *C2
}

func (*Migration2) TableName() string {
	return "migration"
}

type Migration3 struct {
	A      string
	B      string
	D      string
	P1, P2 *C2
}

func (*Migration3) TableName() string {
	return "migration"
}
