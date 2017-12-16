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

func (g *Gen) Translate(name string) string {
	switch name {
	case "AUTO_INCREMENT":
		return "AUTOINCREMENT"
	default:
		return name
	}
}

func (g *Gen) PreProcess(f *load.Field, sqlType *sqltypes.Type) error {
	if f.PrimaryKey {
		if sqlType.Name != sqltypes.Integer {
			return fmt.Errorf("sqlite3 supports autoincrement only for 'INTEGER PRIMARY KEY' columns")
		}
		f.AutoIncrement = true
	}
	return nil
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
