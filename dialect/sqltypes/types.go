package sqltypes

import (
	"regexp"
	"strconv"
)

type Type string

var typeFormat = regexp.MustCompile(`([^(]+)(\((\d+)\))?`)

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

func (t Type) Family() Type {
	m := typeFormat.FindStringSubmatch(string(t))
	return Type(m[1])
}

func (t Type) Size() int {
	m := typeFormat.FindStringSubmatch(string(t))
	if len(m) < 4 {
		return 0
	}
	i, _ := strconv.Atoi(m[3])
	return i
}
