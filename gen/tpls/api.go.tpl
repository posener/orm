package {{.Package}}

import (
	"database/sql"
	"fmt"

    "{{.Type.ImportPath}}"
)

// API is the interface of the ORM object
type API interface {
    Create() *TCreate
    Select() *TSelect
    Insert() *TInsert
    Update() *TUpdate
    Delete() *TDelete
    Insert{{.Type.Name}}(*{{.Type.FullName}}) *TInsert
    Update{{.Type.Name}}(*{{.Type.FullName}}) *TUpdate
}

// Querier is the interface for a SELECT SQL statement
type Querier interface {
    fmt.Stringer
    Query() ([]{{.Type.FullName}}, error)
}

// Execer is the interface for SQL update operations
type Execer interface {
    fmt.Stringer
	Exec() (sql.Result, error)
}

