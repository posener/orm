package {{.Package}}

import (
    "context"

    "{{.Type.ImportPath}}"
)

// API is the interface of the ORM object
type API interface {
    Close() error
    Create() *Create
    Select() *Select
    Insert() *Insert
    Update() *Update
    Delete() *Delete
    Insert{{.Type.Name}}(*{{.Type.FullName}}) *Insert
    Update{{.Type.Name}}(*{{.Type.FullName}}) *Update

    Logger(Logger)
}

// Querier is the interface for a SELECT SQL statement
type Querier interface {
    Query(context.Context) ([]{{.Type.FullName}}, error)
}

// Counter is the interface for a SELECT SQL statement for counting purposes
type Counter interface {
    Count(context.Context) ([]{{.Type.Name}}Count, error)
}
