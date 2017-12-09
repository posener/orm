package load

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/tags"
)

const tagSQLType = "sql"

// Field is a struct that represents type's field
type Field struct {
	ParentType *Type
	// Type is the type of the field
	Type Type
	// Name is the field name
	Name string
	// Embedded means that the field is embedded in a struct
	Embedded bool
	// CustomType can be defined to a column
	CustomType sqltypes.Type
	// PrimaryKey defines a column as a table's primary key
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
	Default    string
	ForeignKey *ForeignKey
}

// ForeignKey is a definition of how a column is a foreign key of another column
// in a referenced table.
type ForeignKey struct {
	Src, Dst *Field
}

func newField(parent *Type, st *types.Struct, i int) (*Field, error) {
	field := st.Field(i)
	if !field.Exported() {
		return nil, nil
	}

	log.Printf("loading field %s", field.Name())

	fieldType, err := New(field.Type().String())
	if err != nil {
		return nil, fmt.Errorf("creating type %s: %s", fieldType, err)
	}

	f := &Field{
		ParentType: parent,
		Name:       field.Name(),
		Type:       *fieldType,
		Embedded:   field.Anonymous(),
	}

	err = f.parseTags(st.Tag(i))

	// if a slice of type 'A', find field in type 'A' that points to the type that the field belongs to.
	switch {
	case f.Type.Slice:
		if f.Type.IsBasic() {
			log.Printf("Ignoring field %s: slice of a basic type is not supported", f.Name)
			return nil, nil
		}
		err = f.findExternalForeignKey()
		if err != nil {
			return nil, err
		}
	case f.IsReference():
		err = f.findForeignKey()
		if err != nil {
			return nil, err
		}
	}

	// set primary key for parent type
	if f.PrimaryKey {
		f.ParentType.PrimaryKey = f
	}

	return f, err
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
			f.CustomType = sqltypes.Type(value)
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
		case "foreign key", "foreign_key", "foreignkey":
			if err := f.setForeignKey(value); err != nil {
				return fmt.Errorf("foreign key definition: %s", err)
			}
		}
	}
	return nil
}

func (f *Field) setForeignKey(name string) error {
	typeName, fieldName := splitForeignKeyTag(name)
	foreignType, err := New(typeName)
	if err != nil {
		return err
	}
	foreignField := foreignType.PrimaryKey
	if fieldName != "" {
		for _, field := range foreignType.Fields {
			if field.Name == fieldName {
				foreignField = field
				break
			}
		}
	}
	if foreignField == nil {
		return fmt.Errorf("no column to reference in foregin table, table should have a primary key, or foreign key definition should incloud foreign column: <type name>#<field name>")
	}
	f.ForeignKey = &ForeignKey{Src: f, Dst: foreignField}
	return nil
}

// findExternalForeignKey looks for foreign key in type that points to this type
// this is useful for slices of one-to-many relation
func (f *Field) findExternalForeignKey() error {
	for _, other := range f.Type.Fields {
		if fk := other.ForeignKey; fk != nil && fk.Dst.ParentType.Table() == f.ParentType.Table() {
			f.ForeignKey = &ForeignKey{Src: fk.Dst, Dst: fk.Src}
			break
		}
	}
	if f.ForeignKey == nil {
		return fmt.Errorf("slice field %s -> %s: did not found foreign key in foreign type %s",
			f.ParentType.ExtNaked(""), f.Name, f.Type.ExtNaked(""))
	}
	return nil
}

// findForeignKey finds the type that this field points to
func (f *Field) findForeignKey() error {
	pk := f.Type.PrimaryKey
	if pk == nil {
		return fmt.Errorf("field %s -> %s: points to type %s which does not have a primary key",
			f.ParentType.ExtNaked(""), f.Name, f.Type.ExtNaked(""))
	}
	f.ForeignKey = &ForeignKey{Src: f, Dst: pk}
	return nil
}

func splitForeignKeyTag(name string) (typeName, fieldName string) {
	i := strings.LastIndex(name, "#")
	if i == -1 {
		return name, ""
	}
	return name[:i], name[i+1:]
}

// Is reference returns true of the field references another row in a table (another object)
func (f *Field) IsReference() bool {
	return !(f.Type.IsBasic() || f.Embedded)
}

// IsSettable returns whether the column could be set
func (f *Field) IsSettable() bool {
	return !(f.PrimaryKey || f.AutoIncrement || f.Type.Slice)
}

// SetType is the type that is used to set this field.
// it is usually the actual type, but in case of reference it is the PK of that type.
func (f *Field) SetType() *Type {
	if f.IsReference() {
		return &f.Type.PrimaryKey.Type
	}
	return &f.Type
}

// Column returns the SQL column name of a field
func (f *Field) Column() string {
	return strings.ToLower(f.Name)
}
