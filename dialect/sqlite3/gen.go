package sqlite3

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/load"
)

// Gen is code generator for sqlite3 dialect
type Gen struct{}

// Name returns the dialect name
func (g *Gen) Name() string {
	return "sqlite3"
}

func (g *Gen) ColumnCreateString(name string, f *load.Field, sqlType *sqltypes.Type) string {
	stmt := []string{fmt.Sprintf("`%s` %s", name, sqlType)}
	if f.NotNull {
		stmt = append(stmt, "NOT NULL")
	}
	if f.Null {
		stmt = append(stmt, "NULL")
	}
	if f.Default != "" {
		stmt = append(stmt, "DEFAULT", f.Default)
	}
	if f.PrimaryKey || f.AutoIncrement {
		stmt = append(stmt, "PRIMARY KEY")
	}
	if f.AutoIncrement {
		if !f.PrimaryKey || sqlType.Name != sqltypes.Integer {
			log.Fatalf("Gen supports autoincrement only for 'INTEGER PRIMARY KEY' columns")
		}
		stmt = append(stmt, "AUTOINCREMENT")
	}
	if f.Unique {
		stmt = append(stmt, " UNIQUE")
	}
	return strings.Join(stmt, " ")
}

func (*Gen) GoTypeToColumnType(t *load.Type) *sqltypes.Type {
	st := new(sqltypes.Type)
	switch typeName := t.Naked.Ext(""); typeName {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		st.Name = sqltypes.Integer
	case "float", "float8", "float16", "float32", "float64":
		st.Name = sqltypes.Float
	case "bool":
		st.Name = sqltypes.Boolean
	case "string":
		st.Name = sqltypes.Text
	case "[]byte":
		st.Name = sqltypes.Blob
	case "time.Time":
		st.Name = sqltypes.TimeStamp
	default:
		log.Fatalf("Unknown column type for %s", typeName)
	}
	return st
}

// ConvertValueCode returns go code for converting value returned from the
// database to the given field.
func (g *Gen) ConvertValueCode(tp *load.Type, field *load.Field, sqlType *sqltypes.Type) string {
	s := tmpltType{
		Type:        tp,
		Field:       field,
		ConvertType: g.convertType(field, sqlType),
	}
	b := bytes.NewBuffer(nil)
	err := tmplt.Execute(b, s)
	if err != nil {
		log.Fatalf("executing sqlite convert value template: %s", err)
	}
	return strings.Trim(b.String(), "\n")
}

type tmpltType struct {
	ConvertType string
	Type        *load.Type
	Field       *load.Field
}

var tmplt = template.Must(template.New("sqlite3").Parse(`
				val, ok := vals[i].({{.ConvertType}})
				if !ok {
					return nil, 0, common.ErrConvert("{{.Field.AccessName}}", i, vals[i], "{{.Field.Type.Ext .Type.Package}}")
				}
				tmp := {{.Field.Type.Naked.Ext .Type.Package}}(val)
				row.{{.Field.AccessName}} = {{if .Field.Type.Pointer -}}&{{end}}tmp
`))

// convertType is the type of the field when returned by sql/driver from database
func (g *Gen) convertType(f *load.Field, sqlType *sqltypes.Type) string {
	switch sqlType.Name {
	case sqltypes.Integer:
		return "int64"
	case sqltypes.Float:
		return "float64"
	case sqltypes.Text, sqltypes.Blob, sqltypes.VarChar:
		return "[]byte"
	case sqltypes.Boolean:
		return "bool"
	default:
		return f.Type.Naked.Ext("")
	}
}
