package load

import "strings"

// Field is a struct that represents type's field
type Field struct {
	// Type is the type of the field
	Type Type
	// Name is the field name
	Name string
	// Embedded means that the field is embedded in a struct
	Embedded bool
	// SQL properties of the field
	SQL SQL
}

// Is reference returns true of the field references another row in a table (another object)
func (f *Field) IsReference() bool {
	return !f.Type.IsBasic() && !f.Embedded
}

// IsSettable returns whether the column could be set
func (f *Field) IsSettable() bool {
	return !f.SQL.PrimaryKey && !f.SQL.AutoIncrement
}

func (f *Field) SetType() string {
	if f.IsReference() {
		return f.Type.PrimaryKey.Type.Name
	}
	return f.Type.ExtName()
}

// Column returns the SQL column name of a field
func (f *Field) Column() string {
	if f.IsReference() {
		return f.Type.Table() + "_id"
	}
	return strings.ToLower(f.Name)
}
