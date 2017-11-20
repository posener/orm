package gen

import (
	"fmt"
	"go/types"
	"log"
	"path/filepath"
	"strings"
)

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

func (t Type) Table() string {
	return strings.ToLower(t.Name)
}

// SQLCreateString returns the SQL CREATE statement for the type
func (t Type) CreateString() string {
	var args = make([]string, 0, len(t.Fields))
	for _, f := range t.Fields {
		args = append(args, f.CreateString())
	}
	return fmt.Sprintf("CREATE TABLE %s ( %s )", t.Table(), strings.Join(args, ", "))
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

// Field is a struct that represents type's field
type Field struct {
	Name       string
	Type       string
	ColumnName string
	// SQL is the SQL properties of the field
	SQL Tags
}

func (f Field) CreateString() string {
	parts := []string{f.ColumnName, f.SQL.Type}
	if f.SQL.NotNull {
		parts = append(parts, "NOT NULL")
	}
	if f.SQL.Default != "" {
		parts = append(parts, "DEFAULT", f.SQL.Default)
	}
	if f.SQL.AutoIncrement {
		parts = append(parts, "AUTO_INCREMENT")
	}
	if f.SQL.Unique {
		parts = append(parts, "UNIQUE")
	}
	if f.SQL.PrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	}
	return strings.Join(parts, " ")
}

func collectFields(s *types.Struct) []Field {
	var f []Field
	for i := 0; i < s.NumFields(); i++ {
		v := s.Field(i)
		if !v.Exported() {
			continue
		}
		varType := v.Type().String()
		tags := ParseTags(s.Tag(i))
		if tags.Type == "" {
			tags.Type = defaultSQLTypes[varType]
		}
		if tags.Type == "" {
			log.Fatalf("Unsupported field type: %s", varType)
		}
		log.Printf("Field '%s(%s)': '%+v'", v.Name(), varType, tags)
		f = append(f, Field{
			Name:       v.Name(),
			Type:       varType,
			ColumnName: strings.ToLower(v.Name()),
			SQL:        tags,
		})
	}
	return f
}
