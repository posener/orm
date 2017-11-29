package load

import (
	"log"
	"strings"

	"go/types"
)

// Type describes the type of the given struct to generate code for
type Type struct {
	GoType
	// Fields is the list of exported fields
	Fields []Field
}

func (t Type) Name() string {
	return t.Type
}

// Field is a struct that represents type's field
type Field struct {
	GoType
	Name string
	SQL  SQL
}

// newType returns a new Type
func newType(name string, st *types.Struct, pkg *types.Package) *Type {
	return &Type{
		GoType: GoType{
			ImportPath: pkg.Path(),
			Type:       name,
		},
		Fields: collectFields(st),
	}
}

// Table is SQL table name of a type
func (t Type) Table() string {
	return strings.ToLower(t.Name())
}

// FieldsImports returns a list of all imports for this type's fields
func (t Type) FieldsImports() []string {
	impsMap := map[string]bool{}
	for _, f := range t.Fields {
		if f.ImportPath != "" {
			impsMap[f.ImportPath] = true
		}
	}
	imps := make([]string, 0, len(impsMap))
	for imp := range impsMap {
		imps = append(imps, imp)
	}
	return imps
}

func collectFields(st *types.Struct) []Field {
	var fields []Field
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		if !field.Exported() {
			continue
		}
		goType := newGoType(field.Type().String())

		switch {
		case field.Anonymous():
			// For embedded fields (aka anonymous) we collect all their fields recursively
			// to the parent fields.
			tp, err := Load(goType)
			if err != nil {
				log.Fatalf("Failed loading type %s: %s", goType, err)
			}
			for _, field := range tp.Fields {
				fields = append(fields, field)
			}

		default:
			sql, err := newSQL(field.Name(), st, i)
			if err != nil {
				log.Fatalf("Failed creating SQL properties for type field %s: %s", goType, err)
			}

			log.Printf("Field '%s(%s)': '%+v'", field.Name(), goType, sql)
			fields = append(fields, Field{
				Name:   field.Name(),
				GoType: goType,
				SQL:    *sql,
			})
		}
	}
	return fields
}
