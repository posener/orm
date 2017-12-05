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

// ConvertValueCode returns go code for converting value returned from the
// database to the given field.
func (g *Gen) ConvertValueCode(field *load.Field) string {
	s := tmpltType{
		Field:             field,
		ConvertType:       g.convertType(field),
		ConvertFuncString: g.convertFuncString(field),
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
	Field             *load.Field
}

var tmplt = template.Must(template.New("sqlite3").Parse(`
				switch val := vals[i].(type) {
				case []byte:
					tmp := {{.ConvertFuncString}}
					row.{{.Field.Name}} = {{if .Field.Type.Pointer}}&{{end}}tmp
				{{- if ne .ConvertType "[]byte" }}
				case {{.ConvertType}}:
					tmp := {{.Field.Type.ExtNaked}}(val)
					row.{{.Field.Name}} = {{if .Field.Type.Pointer -}}&{{end}}tmp
				{{- end }}
				default:
					return nil, fmt.Errorf(errMsg, "{{.Field.Name}}", i, vals[i], vals[i], "[]byte, {{.ConvertType}}")
				}
`))

// convertFuncString is a function for converting the data from SQL to the right type
func (g *Gen) convertFuncString(f *load.Field) string {
	switch tp := strings.TrimLeft(f.SetType(), "*"); tp {
	case "string":
		return "string(val)"
	case "[]byte":
		return "[]byte(val)"
	case "int", "int8", "int16", "int32", "int64":
		return fmt.Sprintf("%s(parseInt(val))", tp)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("%s(parseFloat(val))", tp)
	case "time.Time":
		return fmt.Sprintf("parseTime(val, %d)", g.sqlType(f).Size())
	case "bool":
		return "parseBool(val)"
	default:
		return fmt.Sprintf("%s(val)", tp)
	}
}

func (g *Gen) convertType(f *load.Field) string {
	switch g.sqlType(f).Family() {
	case sqltypes.Integer:
		return "int64"
	case sqltypes.Float:
		return "float64"
	case sqltypes.Text, sqltypes.Blob, sqltypes.VarChar:
		return "[]byte"
	case sqltypes.Boolean:
		return "bool"
	default:
		return f.Type.ExtNaked()
	}
}

func (Gen) sqlType(f *load.Field) sqltypes.Type {
	if custom := f.CustomType; custom != "" {
		return custom
	}
	switch typeName := strings.TrimLeft(f.SetType(), "*"); typeName {
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
