package {{.Package}}

import (
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

// Counter is the interface for a SELECT SQL statement for counting purposes
type Counter interface {
    fmt.Stringer
    Count() ([]{{.Type.Name}}Count, error)
}
