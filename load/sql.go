package load

import (
	"encoding/json"
	"fmt"
	"go/types"
	"strings"

	"github.com/posener/orm/dialect/sqltypes"
	"github.com/posener/orm/tags"
)

const tagSQLType = "sql"

// SQL hold the SQL tags for a field in a struct
type SQL struct {
	// CustomType can be defined to a column
	CustomType sqltypes.Type `json:",omitempty"`
	// PrimaryKey defines a column as a table's primary key
	PrimaryKey bool `json:",omitempty"`
	// NotNull defines that this column value can't be null
	NotNull bool `json:",omitempty"`
	// AutoIncrement defines this column as auto-increment column
	AutoIncrement bool `json:",omitempty"`
	// Unique defines that 2 rows can't have the same value of this column
	Unique bool `json:",omitempty"`
	// Default sets a default value for this column
	Default    string      `json:",omitempty"`
	ForeignKey *ForeignKey `json:",omitempty"`
}

type ForeignKey struct {
	Type  *Type
	Field *Field
}

func newSQL(st *types.Struct, fieldIndex int) (*SQL, error) {
	var sql = new(SQL)
	err := sql.parseTags(st.Tag(fieldIndex))
	return sql, err
}

// parseTags parses tags from a struct tags into a SQL struct.
func (s *SQL) parseTags(tag string) error {
	if tag == "" {
		return nil
	}

	tagsMap := tags.Parse(tag)
	for key, value := range tagsMap[tagSQLType] {
		switch key {
		case "type":
			s.CustomType = sqltypes.Type(value)
		case "primary key", "primary_key", "primarykey":
			s.PrimaryKey = true
		case "not null", "not_null":
			s.NotNull = true
		case "auto increment", "auto_increment", "autoincrement":
			s.AutoIncrement = true
		case "unique":
			s.Unique = true
		case "default":
			s.Default = value
		case "foreign key", "foreign_key", "foreignkey":
			var err error
			s.ForeignKey, err = newForeignKey(value)
			if err != nil {
				return fmt.Errorf("foreign key definition: %s", err)
			}
		}
	}
	return nil
}

func newForeignKey(name string) (*ForeignKey, error) {
	typeName, fieldName := splitForeignKeyTag(name)
	foreignType, err := New(typeName)
	if err != nil {
		return nil, err
	}
	foreignField := foreignType.PrimaryKey
	if fieldName != "" {
		for _, field := range foreignType.Fields {
			if field.Name == fieldName {
				foreignField = &field
				break
			}
		}
	}
	if foreignField == nil {
		return nil, fmt.Errorf("no column to reference in foregin table, table should have a primary key, or foreign key definition should incloud foreign column: <type name>#<field name>")
	}
	return &ForeignKey{Type: foreignType, Field: foreignField}, nil
}

func splitForeignKeyTag(name string) (typeName, fieldName string) {
	i := strings.LastIndex(name, "#")
	if i == -1 {
		return name, ""
	}
	return name[:i], name[i+1:]
}

func (s *SQL) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}
