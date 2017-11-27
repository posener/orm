package common

import (
	"fmt"
	"go/types"
	"log"
	"path/filepath"
	"strings"

	"github.com/posener/orm/load"
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
func NewType(st *load.Struct) Type {
	return Type{
		ImportPath: st.Pkg.Path(),
		Name:       st.Name,
		Fields:     collectFields(st),
	}
}

// Table is SQL table name of a type
func (t Type) Table() string {
	return strings.ToLower(t.Name)
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
	Name string
	Type string
	SQL  SQL
	// ImportPath is a path to add to the import section for this type
	ImportPath string
}

// IsPointerType returns true if field is a pointer
func (f *Field) IsPointerType() bool {
	return len(f.Type) > 0 && f.Type[0] == '*'
}

// NonPointerType returns the non-pointer type of a filed.
// ex, if the type is `*int`, this function will return `int`
func (f *Field) NonPointerType() string {
	if f.IsPointerType() {
		return f.Type[1:]
	}
	return f.Type
}

func collectFields(st *load.Struct) []Field {
	var fields []Field
	for i := 0; i < st.Struct.NumFields(); i++ {
		field := st.Struct.Field(i)
		if !field.Exported() {
			continue
		}
		fieldType := field.Type().String()
		sql, err := newSQL(field.Name(), st.Struct, i)
		if err != nil {
			log.Fatalf("Creating SQL properties for type field %s: %s", st.Name, err)
		}

		log.Printf("Field '%s(%s)': '%+v'", field.Name(), fieldType, sql)
		fields = append(fields, Field{
			Name:       field.Name(),
			Type:       fieldType,
			SQL:        *sql,
			ImportPath: fieldImportPath(fieldType, st.PkgMap),
		})
	}
	return fields
}

func fieldImportPath(typeName string, pkgMap map[string]*types.Package) string {
	i := strings.IndexRune(typeName, '.')
	if i == -1 {
		return ""
	}
	pkg := pkgMap[typeName[:i]]
	if pkg == nil {
		return ""
	}
	return pkg.Path()
}
