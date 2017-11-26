package sqlite3

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/posener/orm/common"
	"github.com/posener/orm/dialect/sqltypes"
)

// Gen is code generator for sqlite3 dialect
type Gen struct {
	Tp common.Type
}

func (g *Gen) Name() string {
	return "sqlite3"
}

func (g *Gen) ConvertValueCode(field *common.Field) string {
	s := tmpltType{
		Field:       field,
		ConvertType: g.convertType(field),
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
	Field       *common.Field
}

var tmplt = template.Must(template.New("sqlite3").Parse(`
				val, ok := vals[i].({{.ConvertType}})
				if !ok {
					return nil, fmt.Errorf(errMsg, "{{.Field.Name}}", i, vals[i], vals[i], "{{.Field.Type}}")
				}
				tmp := {{.Field.NonPointerType}}(val)
				row.{{.Field.Name}} = {{if .Field.IsPointerType -}}&{{end}}tmp
`))

// ConvertType is the type of the field when returned by sql/driver from database
func (g *Gen) convertType(f *common.Field) string {
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
		return f.NonPointerType()
	}
}

func (Gen) sqlType(f *common.Field) sqltypes.Type {
	if f.SQL.CustomType != "" {
		return f.SQL.CustomType
	}
	switch f.NonPointerType() {
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
		return sqltypes.TimeStamp
	default:
		return sqltypes.NA
	}
}
