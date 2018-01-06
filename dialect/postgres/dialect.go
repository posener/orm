package postgres

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/load"
)

// Dialect is code generator for mysql dialect
type Dialect struct{}

// Name returns the dialect name
func (d *Dialect) Name() string {
	return "postgres"
}

func (d *Dialect) Translate(name string) string {
	if name == "AUTO_INCREMENT" {
		return ""
	}
	return name
}

func (d *Dialect) Quote(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

func (Dialect) GoTypeToColumnType(goTypeName string, options []string) *sqltypes.Type {
	st := new(sqltypes.Type)
	if hasAutoIncrement(options) {
		return &sqltypes.Type{Name: "SERIAL"}
	}
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
		st.Name = sqltypes.Text
	case "time.Time":
		st.Name = sqltypes.TimeStamp
		st.Size = 3
	default:
		log.Fatalf("Unknown column type for %s", goTypeName)
	}
	return st
}

func hasAutoIncrement(options []string) bool {
	for _, opt := range options {
		if opt == "AUTO_INCREMENT" {
			return true
		}
	}
	return false
}

// ConvertValueCode returns go code for converting value returned from the
// database to the given field.
func (d *Dialect) ConvertValueCode(field *load.Field) string {
	sqlType := field.CustomType
	if sqlType == nil {
		sqlType = d.GoTypeToColumnType(field.Type.Naked.Ext(""), nil)
	}
	s := tmpltType{
		Field:                  field,
		ConvertType:            d.convertType(field, sqlType),
		ConvertBytesFuncString: d.convertBytesFuncString(field, sqlType),
		ConvertIntFuncString:   d.convertIntFuncString(field, sqlType),
	}
	b := bytes.NewBuffer(nil)
	err := tmplt.Execute(b, s)
	if err != nil {
		log.Fatalf("executing sqlite convert value template: %s", err)
	}
	return strings.Trim(b.String(), "\n")
}

type tmpltType struct {
	ConvertBytesFuncString string
	ConvertIntFuncString   string
	ConvertType            string
	Field                  *load.Field
}

var tmplt = template.Must(template.New("mysql").Parse(`
				switch val := vals[i].(type) {
				case []byte:
					tmp := {{.ConvertBytesFuncString}}
					row.{{.Field.AccessName}} = {{if .Field.Type.Pointer}}&{{end}}tmp
				{{ if ne .ConvertType "[]byte" -}}
				case {{.ConvertType}}:
					tmp := {{.Field.Type.Naked.Ext .Field.ParentType.Package}}(val)
					row.{{.Field.AccessName}} = {{if .Field.Type.Pointer }}&{{end}}tmp
				{{ end -}}
				{{ if and (ne .ConvertType "int64") .ConvertIntFuncString -}}
				case int64:
					tmp := {{.ConvertIntFuncString}}
					row.{{.Field.AccessName}} = {{if .Field.Type.Pointer}}&{{end}}tmp
				{{ end -}}
				default:
					return nil, 0, runtime.ErrConvert("{{.Field.AccessName}}", i, vals[i], "{{.ConvertType}}, []byte, (int64?)")
				}
`))

// convertFuncString is a function for converting the data from SQL to the right type
func (d *Dialect) convertBytesFuncString(f *load.Field, sqlType *sqltypes.Type) string {
	switch tp := f.Type.Naked.Ext(""); tp {
	case "string":
		return "string(val)"
	case "[]byte":
		return "[]byte(val)"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("%s(runtime.ParseInt(val))", tp)
	case "float", "float8", "float16", "float32", "float64":
		return fmt.Sprintf("%s(runtime.ParseFloat(val))", tp)
	case "time.Time":
		return fmt.Sprintf("runtime.ParseTime(val, %d)", sqlType.Size)
	case "bool":
		return "runtime.ParseBool(val)"
	default:
		return fmt.Sprintf("%s(val)", tp)
	}
}

// convertFuncString is a function for converting the data from SQL to the right type
func (d *Dialect) convertIntFuncString(f *load.Field, sqlType *sqltypes.Type) string {
	switch tp := f.Type.Naked.Ext(""); tp {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("%s(val)", tp)
	case "time.Time":
		return fmt.Sprintf("time.Unix(val, 0)")
	case "bool":
		return "val != 0"
	default:
		return ""
	}
}

func (d *Dialect) convertType(f *load.Field, sqlType *sqltypes.Type) string {
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
