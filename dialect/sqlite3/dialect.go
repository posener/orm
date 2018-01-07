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

// Dialect is code generator for sqlite3 dialect
type Dialect struct{}

// Name returns the dialect name
func (d *Dialect) Name() string {
	return "sqlite3"
}

func (d *Dialect) Translate(name string) string {
	// https://sqlite.org/autoinc.html
	if name == "AUTO_INCREMENT" {
		return ""
	}
	return name
}

func (d *Dialect) Quote(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (d *Dialect) ReplaceVars(s string) string {
	return s
}

func (*Dialect) GoTypeToColumnType(goTypeName string, autoIncrement bool) *sqltypes.Type {
	if autoIncrement {
		if !strings.HasPrefix(goTypeName, "int") && !strings.HasPrefix(goTypeName, "uint") {
			log.Panicf("Auto increment for type %v type is not supported", goTypeName)
		}
		return &sqltypes.Type{Name: "integer"}
	}
	st := new(sqltypes.Type)
	switch goTypeName {
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
		st.Name = "integer"
	case "int64", "uint64":
		st.Name = "bigint"
	case "float", "float8", "float16", "float32", "float64":
		st.Name = "real"
	case "bool":
		st.Name = "boolean"
	case "string":
		st.Name = "text"
	case "[]byte":
		st.Name = "blob"
	case "time.Time":
		st.Name = "datetime"
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
		Field:       field,
		ConvertType: d.convertType(field, sqlType),
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
	Field       *load.Field
}

var tmplt = template.Must(template.New("sqlite3").Parse(`
				{{ $convertType := .ConvertType }}
				val, ok := vals[i].({{$convertType}})
				if !ok {
					return nil, 0, dialect.ErrConvert("{{.Field.AccessName}}", i, vals[i], "{{$convertType}}")
				}
				tmp := {{.Field.Type.Naked.Ext .Field.ParentType.Package}}(val)
				row.{{.Field.AccessName}} = {{if .Field.Type.Pointer -}}&{{end}}tmp
`))

// convertType is the type of the field when returned by sql/driver from database
func (d *Dialect) convertType(f *load.Field, sqlType *sqltypes.Type) string {
	switch sqlType.Name {
	case "integer", "bigint":
		return "int64"
	case "real":
		return "float64"
	case "text", "blob", "varchar":
		return "[]byte"
	case "boolean":
		return "bool"
	default:
		return f.Type.Naked.Ext("")
	}
}
