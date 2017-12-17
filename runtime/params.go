package runtime

import (
	"context"
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
