package sqltypes

import "fmt"

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
