package {{.Package}}

import "{{.Type.ImportPath}}"

type API interface {
    Create() *TCreate
    Query() *Select
    Insert() *TInsert
    Update() *TUpdate
    Delete() *TDelete
    Insert{{.Type.Name}}(*{{.Type.FullName}}) *TInsert
    Update{{.Type.Name}}(*{{.Type.FullName}}) *TUpdate
}

