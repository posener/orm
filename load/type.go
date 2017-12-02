package load

import (
	"log"
	"path/filepath"
	"strings"

	"fmt"
	"go/types"
)

var basicTypes = map[string]bool{
	"bool":      true,
	"int":       true,
	"int8":      true,
	"int16":     true,
	"int32":     true,
	"int64":     true,
	"uint":      true,
	"uint8":     true,
	"uint16":    true,
	"uint32":    true,
	"uint64":    true,
	"float32":   true,
	"float64":   true,
	"string":    true,
	"[]byte":    true,
	"time.Time": true,
}

// Type represents a go type attributes by it's name
type Type struct {
	Name string
	// ImportPath is a path to add to the import section for this type
	ImportPath string
	// Fields is the list of exported fields
	Fields     []Field
	PrimaryKey *Field
	Pointer    bool
	Slice      bool
}

// New loads a Type
func New(fullName string) (*Type, error) {
	// []byte is different than any other type since it is allowed slice field
	// which does not actually considered a slice
	if fullName == "[]byte" || fullName == "*[]byte" {
		return &Type{Name: strings.TrimLeft(fullName, "*"), Pointer: pointer(fullName)}, nil
	}
	t := &Type{
		Name:       typeName(fullName),
		ImportPath: importPath(fullName),
		Pointer:    pointer(fullName),
		Slice:      slice(fullName),
	}

	// if type is a basic type, we are done
	if t.IsBasic() {
		return t, nil
	}

	// load the struct data and package information by scanning the go code
	st, pkg, err := t.loadStruct()
	if err != nil {
		return nil, err
	}

	// update the import path to the full package path
	t.ImportPath = pkg.Path()

	// now that we have the type's full name...
	// before loading the fields, check if the type cached in the cache already
	t, cached := cacheGetOrUpdate(t)
	if !cached {
		// the type was not in the cache, we should load all it's fields, which might lead
		// to recursive calls to New function
		err = t.loadFields(st)
	}
	return t, err
}

func (t *Type) String() string {
	if t.ImportPath != "" {
		return t.sliceStr() + t.pointerStr() + t.ImportPath + "." + t.Name
	}
	return t.sliceStr() + t.pointerStr() + t.Name
}

// Table is SQL table name of a type
func (t *Type) Table() string {
	return strings.ToLower(t.Name)
}

// ExtName is the full type of the imported type, as used in a go code
// outside the defining package. For example: "example.Person"
func (t Type) ExtName() string {
	return t.sliceStr() + t.pointerStr() + t.ExtNaked()
}

// ExtNaked is the full type of the imported type in it's non-pointer form,
// as used in a go code outside the defining package.
// For example: "example.Person"
func (t Type) ExtNaked() string {
	if t.Package() != "" {
		return t.Package() + "." + t.Name
	}
	return t.Name
}

// Package is the package name of the type
// for example, type in "github.com/posener/orm/example" has the package
// name: "example"
func (t Type) Package() string {
	_, pkg := filepath.Split(t.ImportPath)
	return pkg
}

func (t *Type) IsBasic() bool {
	return basicTypes[t.ExtNaked()]
}

// Imports returns a list of all imports for this type's fields
func (t *Type) Imports() []string {
	impsMap := map[string]bool{}
	for _, f := range t.Fields {
		if f.Type.ImportPath != "" {
			impsMap[f.Type.ImportPath] = true
		}
	}
	impsMap[t.ImportPath] = true
	imps := make([]string, 0, len(impsMap))
	for imp := range impsMap {
		imps = append(imps, imp)
	}
	return imps
}

// References returns all reference fields
func (t *Type) References() []Field {
	var refs []Field
	for _, field := range t.Fields {
		if field.IsReference() {
			refs = append(refs, field)
		}
	}
	return refs
}

func (t *Type) pointerStr() string {
	if t.Pointer {
		return "*"
	}
	return ""
}

func (t *Type) sliceStr() string {
	if t.Slice {
		return "[]"
	}
	return ""
}

// loadFields iterate over the type's data structure and load all it's fields
// this function might recursively call to the New function
func (t *Type) loadFields(st *types.Struct) error {
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		if !field.Exported() {
			continue
		}

		log.Printf("%s -> loading %s", t.Name, field.Name())

		fieldType, err := New(field.Type().String())
		if err != nil {
			return fmt.Errorf("creating type %s: %s", fieldType, err)
		}
		sql, err := newSQL(st, i)
		if err != nil {
			return fmt.Errorf("creating SQL properties for type field %s: %s", fieldType, err)
		}

		f := Field{
			Name:     field.Name(),
			Type:     *fieldType,
			SQL:      *sql,
			Embedded: field.Anonymous(),
		}

		log.Printf("%s -> loaded %s (%s): %s", t.Name, field.Name(), fieldType, sql)

		switch {
		case f.Embedded:
			// Embedded field (aka anonymous)
			// collect all their fields recursively to the parent fields.
			for _, field := range fieldType.Fields {
				t.Fields = append(t.Fields, field)
			}
		case f.Type.Slice:
			if f.Type.IsBasic() {
				log.Printf("Ignoring field %s: slice of a basic type is not supported", f.Name)
				continue
			}
			foundForeignKey := false
			for _, field := range f.Type.Fields {
				if fk := field.SQL.ForeignKey; fk != nil && fk.Type == t {
					foundForeignKey = true
					break
				}
			}
			if !foundForeignKey {
				return fmt.Errorf("slice field %s -> %s: did not found foreign key in foreign type %s",
					t.ExtNaked(), f.Name, f.Type.ExtNaked())
			}

		default:
			// Basic type field: just add a field
			if sql.PrimaryKey {
				t.PrimaryKey = &f
			}
			t.Fields = append(t.Fields, f)
		}

	}
	return nil
}

// import path returns the import statement of a type
// If a full type name is 'github.com/posener/orm/example.Person', this
// function will return 'github.com/posener/orm/example' which is what you
// would write in the `import` statement.
func importPath(fullName string) string {
	i := strings.LastIndex(fullName, ".")
	if i == -1 {
		return ""
	}
	return strings.TrimLeft(fullName[:i], "*[]")
}

// typeName returns the type string from a full type name.
// If a full type name is 'github.com/posener/orm/example.Person', this
// function will return 'Person' which is how you would use this
// struct in a go file
func typeName(fullName string) string {
	i := strings.LastIndex(fullName, ".")
	return strings.TrimLeft(fullName[i+1:], "*[]")
}

func pointer(typeName string) bool {
	return strings.HasPrefix(strings.TrimPrefix(typeName, "[]"), "*")
}

func slice(typeName string) bool {
	return strings.HasPrefix(typeName, "[]")
}
