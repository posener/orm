package tests

//go:generate ../orm -type Migration0,Migration1

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
