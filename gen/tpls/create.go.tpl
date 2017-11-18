package {{.Package}}

import (
	"database/sql"
	"log"
)

{{$n := len .Type.Fields}}

// Create creates a table for {{.Type.Name}}
func Create(db *sql.DB) error {
    // Create statement has a line for each variable with it's name and it's type.
    // The if statement is to remove the comma from the last line.
	stmt := `CREATE TABLE {{.Table}} (
	{{- range $i, $f := .Type.Fields }}
		{{$f.ColumnName}} {{$f.ColumnSQLType}}{{if lt (plus1 $i) $n}},{{end}}
    {{- end }}
	)`
	log.Printf("Create: '%v'", stmt)
	_, err := db.Exec(stmt)
	return err
}
