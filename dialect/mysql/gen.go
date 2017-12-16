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

func (g *Gen) Translate(name string) string {
	return name
}

func (Gen) GoTypeToColumnType(goTypeName string) *sqltypes.Type {
	st := new(sqltypes.Type)
	switch goTypeName {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		st.Name = sqltypes.Integer
	case "float", "float8", "float16", "float32", "float64":
		st.Name = sqltypes.Float
	case "bool":
		st.Name = sqltypes.Boolean
	case "string":
		st.Name = sqltypes.VarChar
		st.Size = 255
	case "[]byte":
		st.Name = sqltypes.Blob
	case "time.Time":
		st.Name = sqltypes.DateTime
		st.Size = 3
	default:
		log.Fatalf("Unknown column type for %s", goTypeName)
	}
	return st
}

// ConvertValueCode returns go code for converting value returned from the
// database to the given field.
func (g *Gen) ConvertValueCode(field *load.Field, sqlType *sqltypes.Type) string {
	s := tmpltType{
		Field:             field,
		ConvertType:       g.convertType(field, sqlType),
		ConvertFuncString: g.convertFuncString(field, sqlType),
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

var tmplt = template.Must(template.New("mysql").Parse(`
				switch val := vals[i].(type) {
				case []byte:
					tmp := {{.ConvertFuncString}}
					row.{{.Field.AccessName}} = {{if .Field.Type.Pointer}}&{{end}}tmp
				{{- if ne .ConvertType "[]byte" }}
				case {{.ConvertType}}:
					tmp := {{.Field.Type.Naked.Ext .Field.ParentType.Package}}(val)
					row.{{.Field.AccessName}} = {{if .Field.Type.Pointer -}}&{{end}}tmp
				{{- end }}
				default:
					return nil, 0, runtime.ErrConvert("{{.Field.AccessName}}", i, vals[i], "[]byte, {{.ConvertType}}")
				}
`))

// convertFuncString is a function for converting the data from SQL to the right type
func (g *Gen) convertFuncString(f *load.Field, sqlType *sqltypes.Type) string {
	switch tp := f.Type.Naked.Ext(""); tp {
	case "string":
		return "string(val)"
	case "[]byte":
		return "[]byte(val)"
	case "int", "int8", "int16", "int32", "int64":
		return fmt.Sprintf("%s(runtime.ParseInt(val))", tp)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("%s(runtime.ParseFloat(val))", tp)
	case "time.Time":
		return fmt.Sprintf("runtime.ParseTime(val, %d)", sqlType.Size)
	case "bool":
		return "runtime.ParseBool(val)"
	default:
		return fmt.Sprintf("%s(val)", tp)
	}
}

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
