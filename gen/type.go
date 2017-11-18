package gen

import (
	"fmt"
	"go/types"
	"log"
	"path/filepath"
	"strings"
)

const tagSQLType = "sql"

var defaultSQLTypes = map[string]string{
	"string": "VARCHAR(255)",
	"int":    "INT",
	"bool":   "BOOLEAN",
	"float":  "REAL",
}

// Type describes the type of the given struct to generate code for
type Type struct {
	// ImportPath is the import path of the package of the given struct.
	// for example: "github.com/posener/orm/example"
	ImportPath string
	// Name is the type name, for example, the name of the Person type in
	// the example package is "Person"
	Name string
	// Fields is the list of exported fields
	Fields []Field
}

// NewType returns a new Type
func NewType(name string, pkg *types.Package, st *types.Struct) Type {
	return Type{
		ImportPath: pkg.Path(),
		Name:       name,
		Fields:     collectFields(st),
	}
}

// FullName is the full type of the imported type, as used in a go code
// outside the defining package. For example: "example.Person"
func (t Type) FullName() string {
	return fmt.Sprintf("%s.%s", t.Package(), t.Name)
}

// Package is the package name of the type
// for example, type in "github.com/posener/orm/example" has the package
// name: "example"
func (t Type) Package() string {
	_, pkg := filepath.Split(t.ImportPath)
	return pkg
}

type Field struct {
	Name       string
	ColumnName string
	Type       string
	// ColumnSQLType is the SQL type of the field
	ColumnSQLType string
}

func collectFields(s *types.Struct) []Field {
	var f []Field
	for i := 0; i < s.NumFields(); i++ {
		v := s.Field(i)
		if !v.Exported() {
			continue
		}
		varType := v.Type().String()
		sqlType := customSQLTypeFromTag(varType)
		if sqlType == "" {
			sqlType = defaultSQLTypes[varType]
		}
		if sqlType == "" {
			log.Fatalf("Unsupported field type: %s", varType)
		}
		s.Tag(i)
		f = append(f, Field{
			Name:          v.Name(),
			ColumnName:    strings.ToLower(v.Name()),
			Type:          varType,
			ColumnSQLType: sqlType,
		})
	}
	return f
}

func customSQLTypeFromTag(tag string) string {
	parts := strings.Fields(tag)
	for _, part := range parts {
		keyVal := strings.Split(part, ":")
		if len(keyVal) != 2 {
			continue
		}
		key := keyVal[0]
		val := keyVal[1]
		if key != tagSQLType {
			continue
		}
		return strings.Trim(val, `"`)
	}
	return ""
}
