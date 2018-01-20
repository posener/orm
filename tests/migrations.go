package tests

//go:generate ../orm -type Migration0,Migration1,Migration2,Migration3

// Migration0 is for testing migrations
type Migration0 struct {
	A string
}

// TableName returns the table name of Migration0 type
func (*Migration0) TableName() string {
	return "migration"
}

// Migration1 is for testing migrations
type Migration1 struct {
	A string
	B string
}

// TableName returns the table name of Migration1 type
func (*Migration1) TableName() string {
	return "migration"
}

// Migration2 is for testing migrations
type Migration2 struct {
	A  string
	B  string
	D  string
	P1 *C2
}

// TableName returns the table name of Migration2 type
func (*Migration2) TableName() string {
	return "migration"
}

// Migration3 is for testing migrations
type Migration3 struct {
	A      string
	B      string
	D      string
	P1, P2 *C2
}

// TableName returns the table name of Migration3 type
func (*Migration3) TableName() string {
	return "migration"
}
