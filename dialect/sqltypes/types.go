package sqltypes

import (
	"regexp"
	"strconv"
)

var typeFormat = regexp.MustCompile(`([^(]+)(\((\d+)\))?`)

// Type represents an SQL column type
type Type string

// List of SQL types
const (
	NA        Type = ""
	Integer   Type = "INTEGER"
	Float     Type = "FLOAT"
	Boolean   Type = "BOOLEAN"
	Text      Type = "TEXT"
	Blob      Type = "BLOB"
	TimeStamp Type = "TIMESTAMP"
	DateTime  Type = "DATETIME"
	VarChar   Type = "VARCHAR"
)

// Family returns the family of a type
// ex. the family of VARCHAR(10) is VARCHAR, the family of INT is INT.
func (t Type) Family() Type {
	m := typeFormat.FindStringSubmatch(string(t))
	return Type(m[1])
}

// Size returns the size of a column, if that is defined.
// ex. the size of VARCHAR(10) is 10, the size of INT is 0 since it is not defined.
func (t Type) Size() int {
	m := typeFormat.FindStringSubmatch(string(t))
	if len(m) < 4 {
		return 0
	}
	i, _ := strconv.Atoi(m[3])
	return i
}
