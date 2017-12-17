package runtime

import (
	"context"
	"fmt"
	"strings"
)

// CreateParams holds parameters for an SQL CREATE statement
type CreateParams struct {
	Table          string
	MarshaledTable string
	// IfNotExists determines to create the table only if it does not exists
	IfNotExists bool
	// AutoMigrate perform auto-migration of table scheme
	AutoMigrate bool
	// Ctx is a context parameter for the query
	// even though it is not recommended to store context in a struct, here the struct
	// actually represents an arguments list, passed to a function.
	Ctx context.Context
}

type TableProperties struct {
	Columns     map[string]string
	ForeignKeys map[string]string
	PrimaryKeys []string
}

func (tp TableProperties) String() string {
	stmts := make([]string, 0, len(tp.Columns)+len(tp.ForeignKeys))
	for _, col := range tp.Columns {
		stmts = append(stmts, col)
	}
	if len(tp.PrimaryKeys) > 0 {
		stmts = append(stmts, fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(tp.PrimaryKeys, ", ")))
	}
	for _, fk := range tp.ForeignKeys {
		stmts = append(stmts, fk)
	}
	return strings.Join(stmts, ", ")
}

// InsertParams holds parameters for an SQL INSERT statement
type InsertParams struct {
	Table string
	// Assignments are values to store in the new row
	Assignments Assignments
	// Ctx is a context parameter for the query
	// even though it is not recommended to store context in a struct, here the struct
	// actually represents an arguments list, passed to a function.
	Ctx context.Context
}

// SelectParams holds parameters for an SQL SELECT statement
type SelectParams struct {
	Table   string
	Columns Selector
	Where   Where
	Groups  Groups
	Orders  Orders
	Page    Page
	// Ctx is a context parameter for the query
	// even though it is not recommended to store context in a struct, here the struct
	// actually represents an arguments list, passed to a function.
	Ctx context.Context
}

// DeleteParams holds parameters for an SQL DELETE statement
type DeleteParams struct {
	Table string
	Where Where
	// Ctx is a context parameter for the query
	// even though it is not recommended to store context in a struct, here the struct
	// actually represents an arguments list, passed to a function.
	Ctx context.Context
}

// UpdateParams holds parameters for an SQL UPDATE statement
type UpdateParams struct {
	Table       string
	Assignments Assignments
	Where       Where
	// Ctx is a context parameter for the query
	// even though it is not recommended to store context in a struct, here the struct
	// actually represents an arguments list, passed to a function.
	Ctx context.Context
}
