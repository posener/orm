package postgres

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"regexp"

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

var qMarkReplacer = regexp.MustCompile(`(\?)`)

func (d *Dialect) ReplaceVars(s string) string {
	i := 0
	return qMarkReplacer.ReplaceAllStringFunc(s, func(string) string {
		i++
		return fmt.Sprintf("$%d", i)
	})
}

func (Dialect) GoTypeToColumnType(goTypeName string, autoIncrement bool) *sqltypes.Type {
	st := new(sqltypes.Type)
	if autoIncrement {
		switch goTypeName {
		case "int", "int8", "int16", "int32", "uint", "uint8", "uint16":
			return &sqltypes.Type{Name: "serial"}
		case "int64", "uint32", "uint64":
			return &sqltypes.Type{Name: "bigserial"}
		default:
			log.Fatalf("Auto increment on type %s is not supported", goTypeName)
		}
	}
	switch goTypeName {
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16":
		st.Name = "integer"
	case "int64", "uint32", "uint64":
		st.Name = "bigint"
	case "float32", "float64":
		st.Name = "numeric"
	case "bool":
		st.Name = "boolean"
	case "string":
		st.Name = "varchar"
		st.Size = 255
	case "[]byte":
		st.Name = "text"
	case "time.Time":
		st.Name = "timestamp with time zone"
	default:
		log.Fatalf("Unknown column type for %s", goTypeName)
	}
	return st
}

// ConvertValueCode returns go code for converting value returned from the
// database to the given field.
func (d *Dialect) ConvertValueCode(field *load.Field) string {
	sqlType := field.CustomType
	if sqlType == nil {
		sqlType = d.GoTypeToColumnType(field.Type.Naked.Ext(""), false)
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

var tmplt = template.Must(template.New("postgres").Parse(`
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
	case "float32", "float64":
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
	case "integer", "bigint":
		return "int64"
	case "numeric":
		return "float64"
	case "varchar":
		return "string"
	case "text":
		return "[]byte"
	case "boolean":
		return "bool"
	default:
		return f.Type.Naked.Ext("")
	}
}
