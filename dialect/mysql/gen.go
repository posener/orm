package mysql

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/load"
)

// Gen is code generator for mysql dialect
type Gen struct{}

// Name returns the dialect name
func (g *Gen) Name() string {
	return "mysql"
}

func (g *Gen) ColumnCreateString(f *load.Field, sqlType sqltypes.Type) string {
	stmt := []string{fmt.Sprintf("`%s` %s", f.Column(), sqlType)}
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
		stmt = append(stmt, "AUTO_INCREMENT")
	}
	if f.Unique {
		stmt = append(stmt, " UNIQUE")
	}
	return strings.Join(stmt, " ")

}

func (Gen) GoTypeToColumnType(t *load.Type) sqltypes.Type {
	switch typeName := t.ExtNaked(""); typeName {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return sqltypes.Integer
	case "float", "float8", "float16", "float32", "float64":
		return sqltypes.Float
	case "bool":
		return sqltypes.Boolean
	case "string":
		return sqltypes.Text
	case "[]byte":
		return sqltypes.Blob
	case "time.Time":
		return sqltypes.DateTime + "(3)"
	default:
		log.Fatalf("Unknown column type for %s", typeName)
		return sqltypes.NA
	}
}

// ConvertValueCode returns go code for converting value returned from the
// database to the given field.
func (g *Gen) ConvertValueCode(tp *load.Type, field *load.Field, sqlType sqltypes.Type) string {
	s := tmpltType{
		Type:              tp,
		Field:             field,
		ConvertType:       g.convertType(field, sqlType),
		ConvertFuncString: g.convertFuncString(tp, field, sqlType),
	}
	b := bytes.NewBuffer(nil)
	err := tmplt.Execute(b, s)
	if err != nil {
		log.Fatalf("executing sqlite convert value template: %s", err)
	}
	return strings.Trim(b.String(), "\n")
}

type tmpltType struct {
	ConvertFuncString string
	ConvertType       string
	Type              *load.Type
	Field             *load.Field
}

var tmplt = template.Must(template.New("mysql").Parse(`
				switch val := vals[i].(type) {
				case []byte:
					tmp := {{.ConvertFuncString}}
					row.{{.Field.Name}} = {{if .Field.Type.Pointer}}&{{end}}tmp
				{{- if ne .ConvertType "[]byte" }}
				case {{.ConvertType}}:
					tmp := {{.Field.Type.ExtNaked .Type.Package}}(val)
					row.{{.Field.Name}} = {{if .Field.Type.Pointer -}}&{{end}}tmp
				{{- end }}
				default:
					return nil, common.ErrConvert("{{.Field.Name}}", i, vals[i], "[]byte, {{.ConvertType}}")
				}
`))

// convertFuncString is a function for converting the data from SQL to the right type
func (g *Gen) convertFuncString(t *load.Type, f *load.Field, sqlType sqltypes.Type) string {
	switch tp := f.SetType().ExtNaked(""); tp {
	case "string":
		return "string(val)"
	case "[]byte":
		return "[]byte(val)"
	case "int", "int8", "int16", "int32", "int64":
		return fmt.Sprintf("%s(common.ParseInt(val))", tp)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("%s(common.ParseFloat(val))", tp)
	case "time.Time":
		return fmt.Sprintf("common.ParseTime(val, %d)", sqlType.Size())
	case "bool":
		return "common.ParseBool(val)"
	default:
		return fmt.Sprintf("%s(val)", tp)
	}
}

func (g *Gen) convertType(f *load.Field, sqlType sqltypes.Type) string {
	switch sqlType.Family() {
	case sqltypes.Integer:
		return "int64"
	case sqltypes.Float:
		return "float64"
	case sqltypes.Text, sqltypes.Blob, sqltypes.VarChar:
		return "[]byte"
	case sqltypes.Boolean:
		return "bool"
	default:
		return f.Type.ExtNaked("")
	}
}
