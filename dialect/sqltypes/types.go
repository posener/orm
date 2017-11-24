package sqltypes

import (
	"fmt"
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
)

func VarChar(size int) Type {
	return Type(fmt.Sprintf("VARCHAR(%d)", size))
}

// Family returns the family of a type
// for example, for VARCHAR(???) the family is VARCHAR
// for INTEGER, the family is INTEGER
func (t Type) Family() string {
	s := string(t)
	prefixEnds := strings.Index(s, "(")
	if prefixEnds == -1 {
		prefixEnds = len(s)
	}
	return s[:prefixEnds]
}
