package sqltypes

import (
	"strings"
)

type Type string

const (
	NA        Type = ""
	Integer   Type = "INTEGER"
	Float     Type = "FLOAT"
	Boolean   Type = "BOOLEAN"
	Text      Type = "TEXT"
	Blob      Type = "BLOB"
	TimeStamp Type = "TIMESTAMP"
	VarChar   Type = "VARCHAR"
)

func Family(t Type) Type {
	prefixEnds := strings.Index(string(t), "(")
	if prefixEnds == -1 {
		return t
	}
	return Type(string(t)[:prefixEnds])
}
