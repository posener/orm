package runtime

import (
	"context"
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
	Table string
	// Columns are the selected columns for the query
	Columns map[string]bool
	// Joins are recursive join arguments for the query
	Joins []JoinParams
	// Count indicates weather a COUNT(*) column should be added to the query
	Count bool
	// Where is the WHERE part of the query
	Where Where
	// Groups are is the GROUP BY conditions for the query
	Groups Groups
	// Orders are the ORDER BY conditions for the query
	Orders Orders
	// Page is for query pagination
	Page Page
	// OrderedColumns store all the columns of a table
	// they are defined in a specific order so the parsing of returned values from
	// the SQL query will be easy.
	OrderedColumns []string
	// Ctx is a context parameter for the query
	// even though it is not recommended to store context in a struct, here the struct
	// actually represents an arguments list, passed to a function.
	Ctx context.Context
}

// SelectAll indicates if the SelectParams should select all columns in the query string
func (p *SelectParams) SelectAll() bool {
	return len(p.Columns) == 0 && !p.Count
}

// SelectedColumns returns a list of selected columns for a query
func (p *SelectParams) SelectedColumns() []string {
	// if select all, we insert all the columns into the params, in order to define a specific
	// order in the returned values from the SQL query.
	if p.SelectAll() {
		return p.OrderedColumns
	}
	cols := make([]string, 0, len(p.Columns))
	for _, col := range p.OrderedColumns {
		if p.Columns[col] {
			cols = append(cols, col)
		}
	}
	return cols
}

// JoinParams are parameters to perform a join operation
// Field SelectParams is used to perform select operations on the join struct's field.
// Pairings describe the relation between the join's fields
type JoinParams struct {
	SelectParams
	Pairings []Pairing
}

// TableName creates a table name for a join operation
// this is useful in case several fields referencing the same table
func (j *JoinParams) TableName(parentTable string) string {
	parts := make([]string, 0, len(j.Pairings)+1)
	parts = append(parts, parentTable)
	for _, pairing := range j.Pairings {
		parts = append(parts, pairing.Column)
	}
	return strings.Join(parts, "_")
}

// Pairing describe a join relation
type Pairing struct {
	// Column is the column in the current table for the JOIN statement
	Column string
	// JoinedColumn is the column in the referenced table for the JOIN statement
	JoinedColumn string
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

// DropParams holds parameters for an SQL DROP statement
type DropParams struct {
	Table    string
	IfExists bool
	// Ctx is a context parameter for the query
	// even though it is not recommended to store context in a struct, here the struct
	// actually represents an arguments list, passed to a function.
	Ctx context.Context
}
