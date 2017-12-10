package load

import (
	"fmt"
	"go/types"
	"log"
	"path/filepath"
	"strings"
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
	*Naked
	Pointer bool
	Slice   bool
}

// Naked is a type definition without it's specific usage information
type Naked struct {
	Name string
	// ImportPath is a path to add to the import section for this type
	ImportPath string
	// Fields is the list of exported fields
	Fields     []*Field
	PrimaryKey *Field
}

// New loads a Type
func New(fullName string) (*Type, error) {
	// []byte is different than any other type since it is allowed slice field
	// which does not actually considered a slice
	if fullName == "[]byte" || fullName == "*[]byte" {
		return &Type{
			Naked:   &Naked{Name: strings.TrimLeft(fullName, "*")},
			Pointer: pointer(fullName),
		}, nil
	}
	t := &Type{
		Naked: &Naked{
			Name:       typeName(fullName),
			ImportPath: importPath(fullName),
		},
		Pointer: pointer(fullName),
		Slice:   slice(fullName),
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
	var cached bool
	t.Naked, cached = cacheGetOrUpdate(t.Naked)
	if cached {
		return t, err
	}
	// the type was not in the cache, we should load all it's fields, which might lead
	// to recursive calls to New function
	err = t.loadFields(st)
	if err != nil {
		return nil, err
	}

	log.Printf("%s: creating graph", t.Naked)
	for _, f := range t.Fields {
		err := f.SetRefType()
		if err != nil {
			return nil, fmt.Errorf("set relation for field '%s': %s", f, err)
		}
		if f.RefType != nil {
			log.Printf("Edge: %s -> %s", f, f.RefType)
		}
	}
	return t, nil
}

func (t *Type) SetRelations() error {
	log.Printf("%s: calculating foreign keys", t)
	for _, f := range t.Fields {
		err := f.SetForeignKey()
		if err != nil {
			return fmt.Errorf("set foreign key for field '%s': %s", f, err)
		}
	}
	return nil
}

// loadFields iterate over the type's data structure and load all it's fields
// this function might recursively call to the New function
func (t *Naked) loadFields(st *types.Struct) error {
	for i := 0; i < st.NumFields(); i++ {
		field, err := newField(t, st, i)
		if err != nil {
			return err
		}
		switch {
		case field == nil:
		case field.Embedded:
			// Embedded field (aka anonymous)
			// collect all their fields recursively to the parent fields.
			for _, field := range field.Type.Fields {
				t.Fields = append(t.Fields, field)
			}
		default:
			// Basic type field: just add a field
			t.Fields = append(t.Fields, field)
		}
	}
	return nil
}

func (t *Naked) String() string {
	if t.ImportPath != "" {
		return t.ImportPath + "." + t.Name
	}
	return t.Name
}

func (t *Type) String() string {
	return t.sliceStr() + t.pointerStr() + t.Naked.String()
}

// Table is SQL table name of a type
func (t *Naked) Table() string {
	return strings.ToLower(t.Name)
}

// Ext return the type representation depending on the given package,
// if it is the same package as the type's, it will return only it's
// name. Otherwise, it will return the full "package.Name" semantic
func (t *Type) Ext(curPkg string) string {
	return t.sliceStr() + t.pointerStr() + t.Naked.Ext(curPkg)
}

// Ext return the type representation depending on the given package,
// if it is the same package as the type's, it will return only it's
// name. Otherwise, it will return the full "package.Name" semantic.
func (t *Naked) Ext(curPkg string) string {
	if t.Package() != "" && t.Package() != curPkg {
		return t.Package() + "." + t.Name
	}
	return t.Name
}

// Package is the package name of the type
// for example, type in "github.com/posener/orm/example" has the package
// name: "example"
func (t *Naked) Package() string {
	_, pkg := filepath.Split(t.ImportPath)
	return pkg
}

func (t *Type) IsBasic() bool {
	return basicTypes[t.Naked.Ext("")]
}

// Imports returns a list of all imports for this type's fields
func (t *Type) Imports() []string {
	impsMap := map[string]bool{}
	for _, f := range t.Fields {
		if f.Type.ImportPath != "" && f.Type.ImportPath != t.ImportPath {
			impsMap[f.Type.ImportPath] = true
		}
	}
	imps := make([]string, 0, len(impsMap))
	for imp := range impsMap {
		imps = append(imps, imp)
	}
	return imps
}

// References returns all reference fields
func (t *Type) References() []*Field {
	var refs []*Field
	for _, field := range t.Fields {
		if field.IsReference() {
			refs = append(refs, field)
		}
	}
	return refs
}

func (t *Type) HasOneToManyRelation() bool {
	for _, field := range t.Fields {
		if field.Type.Slice {
			return true
		}
	}
	return false
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
