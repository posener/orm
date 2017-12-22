package load

import (
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/tags"
)

const tagSQLType = "sql"

// Field is a struct that represents type's field
type Field struct {
	ParentType *Naked
	// Type is the type of the field
	Type Type
	// AccessName is the name of the field from the root struct
	// so if a field is embedded, StructA->StructB->Field, the field AccessName, from StructA
	// will be StructB.Field
	AccessName string
	// Embedded means that the field is embedded in a struct
	Embedded bool
	// CustomType can be defined to a column
	CustomType *sqltypes.Type
	// PrimaryKeys defines a column as a table's primary key
	PrimaryKey bool
	// NotNull defines that this column value can't be null
	NotNull bool
	// Null defines that this column value can be null
	Null bool
	// AutoIncrement defines this column as auto-increment column
	AutoIncrement bool
	// Unique defines that 2 rows can't have the same value of this column
	Unique bool
	// Default sets a default value for this column
	Default string
	// ReferencingFieldName is to used when type A has a one-to-many relationship
	// to type B, and type B has more than one fields that reference type A.
	// in this case, the one-to-many field in A, for example 'Bs []B' should add
	// a tag with the name of the field in B that referencing it. So if type
	// B hash 'A1, A2 A', type B should add tag `sql:"referencing field:A1"`.
	ReferencingFieldName string
}

// ForeignKey is a definition of how a column is a foreign key of another column
// in a referenced table.
type ForeignKey struct {
	Src, Dst *Field
}

func newField(parent *Naked, i int) (*Field, error) {
	stField := parent.st.Field(i)
	if !stField.Exported() {
		return nil, nil
	}

	log.Printf("loading field %s", stField.Name())

	fieldType, err := New(stField.Type().String())
	if err != nil {
		return nil, fmt.Errorf("creating type %s: %s", stField.Type().String(), err)
	}

	f := &Field{
		ParentType: parent,
		AccessName: stField.Name(),
		Type:       *fieldType,
		Embedded:   stField.Anonymous(),
	}

	// ignore slice of basic type - not supported
	if f.Type.Slice && f.Type.IsBasic() {
		log.Printf("Ignoring field %s: slice of a basic type is not supported", f.Name())
		return nil, nil
	}

	err = f.parseTags(parent.st.Tag(i))
	if err != nil {
		return nil, fmt.Errorf("%s: parse tags: %s", f, err)
	}

	// set primary key for parent type
	if f.PrimaryKey || f.Unique {
		log.Printf("Field %s: set as primary key", f)
		f.ParentType.PrimaryKeys = append(f.ParentType.PrimaryKeys, f)
	}

	return f, err
}

// Name is the field name
// If the field is embedded withing another type, for example StructA->StructB->Field, to
// distinct between a field with name 'Field' in StructA, the name will be 'StructBField'
func (f *Field) Name() string {
	return strings.Replace(f.AccessName, ".", "", -1)
}

// parseTags parses tags from a struct tags into a SQL struct.
func (f *Field) parseTags(tag string) error {
	if tag == "" {
		return nil
	}

	tagsMap := tags.Parse(tag)
	for key, value := range tagsMap[tagSQLType] {
		switch key {
		case "type":
			var err error
			f.CustomType, err = sqltypes.New(value)
			if err != nil {
				return fmt.Errorf("parsing type %s: %s", value, err)
			}
		case "primary key", "primary_key", "primarykey":
			f.PrimaryKey = true
		case "not null", "not_null":
			f.NotNull = true
		case "null":
			f.Null = true
		case "auto increment", "auto_increment", "autoincrement":
			f.AutoIncrement = true
		case "unique":
			f.Unique = true
		case "default":
			f.Default = value
		case "referencing field", "referencing_field":
			f.ReferencingFieldName = value
		}
	}
	return nil
}

// IsReference returns true of the field references another row in a table (another object)
func (f *Field) IsReference() bool {
	return !(f.Type.IsBasic() || f.Embedded)
}

// IsForwardReference returns true for a type that references another type
func (f *Field) IsForwardReference() bool {
	return f.IsReference() && !f.Type.Slice
}

// IsReversedReference returns true for a type that is referenced by other types
func (f *Field) IsReversedReference() bool {
	return f.IsReference() && f.Type.Slice
}

// IsSettable returns whether the column could be set
func (f *Field) IsSettable() bool {
	return !(f.AutoIncrement || f.Type.Slice)
}

// Columns returns the SQL column name of a field
func (f *Field) Columns() []SQLColumn {
	if f.IsForwardReference() {
		if f.CustomType != nil {
			log.Fatalf("filed %s is reference, can't have a custom type", f)
		}
		cols := make([]SQLColumn, 0, len(f.Type.PrimaryKeys))
		for _, pk := range f.Type.PrimaryKeys {
			cols = append(cols, SQLColumn{
				Name:    strings.ToLower(fmt.Sprintf("fk_%s_%s", f.Name(), pk.Name())),
				SetType: &pk.Type,
			})
		}
		return cols
	}
	return []SQLColumn{f.column()}
}

// Column returns the SQL column name of a field
func (f *Field) Column() SQLColumn {
	if f.IsReference() {
		log.Panic("Column should not be called on a reference field, use columns")
	}
	return f.column()
}
func (f *Field) column() SQLColumn {
	return SQLColumn{
		Name:       strings.ToLower(f.Name()),
		SetType:    &f.Type,
		CustomType: f.CustomType,
	}
}

// SQLColumn describe a column in an SQL table
type SQLColumn struct {
	// Name is the column name
	Name string
	// SetTypes is the type that is used to set a field that reference this column
	SetType *Type
	// CustomType is a custom SQL type that can be defined by the user
	CustomType *sqltypes.Type
}

func (f *Field) String() string {
	return fmt.Sprintf("%s#%s", f.ParentType.Ext(""), f.Name())
}

// PrivateName returns the name of the field with the first letter as a lower case
// so it could be used as a private variable name
func (f *Field) PrivateName() string {
	name := f.Name()
	return strings.ToLower(name[:1]) + name[1:]
}
