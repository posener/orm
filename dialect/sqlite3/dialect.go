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
		return &sqltypes.Type{Name: "INTEGER"}
	}
	st := new(sqltypes.Type)
	switch goTypeName {
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
		st.Name = "INTEGER"
	case "int64", "uint64":
		st.Name = "BIGINT"
	case "float", "float8", "float16", "float32", "float64":
		st.Name = "REAL"
	case "bool":
		st.Name = "BOOL"
	case "string":
		st.Name = "TEXT"
	case "[]byte":
		st.Name = "BLOB"
	case "time.Time":
		st.Name = "DATETIME"
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
				val, ok := vals[i].({{.ConvertType}})
				if !ok {
					return nil, 0, runtime.ErrConvert("{{.Field.AccessName}}", i, vals[i], "{{.Field.Type.Ext .Field.ParentType.Package}}")
				}
				tmp := {{.Field.Type.Naked.Ext .Field.ParentType.Package}}(val)
				row.{{.Field.AccessName}} = {{if .Field.Type.Pointer -}}&{{end}}tmp
`))

// convertType is the type of the field when returned by sql/driver from database
func (d *Dialect) convertType(f *load.Field, sqlType *sqltypes.Type) string {
	switch sqlType.Name {
	case "INTEGER", "BIGINT":
		return "int64"
	case "REAL":
		return "float64"
	case "TEXT", "BLOB":
		return "[]byte"
	case "BOOL":
		return "bool"
	default:
		return f.Type.Naked.Ext("")
	}
}
