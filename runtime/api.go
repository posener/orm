package runtime

import "github.com/posener/orm"

// TableNamer is an interface for a model to change it's table name
// a struct that implements this interface will set it's SQL table name by the return
// value from TableName function.
type TableNamer interface {
	TableName() string
}

// StatementArger is interface for queries.
// The statement and the args are given to the SQL query.
type StatementArger interface {
	Statement(string) string
	Args() []interface{}
}

// Order is a struct that is used for ORDER BY operations.
// It holds the column name and the order direction.
type Order struct {
	Column string
	Dir    orm.OrderDir
}

// Orders is a list of Order
type Orders []Order

// Add adds an order of a column and direction
func (g *Orders) Add(name string, dir orm.OrderDir) {
	*g = append(*g, Order{Column: name, Dir: dir})
}

// Group is a struct that is used for GROUP BY operations
type Group struct {
	Column string
}

// Groups is a list of Group
type Groups []Group

// Add adds a column to a group
func (g *Groups) Add(name string) {
	*g = append(*g, Group{Column: name})
}

// Page is a struct that is used for LIMIT/OFFSET operations
type Page struct {
	Limit  int64
	Offset int64
}

// Assignment is a struct that is used for INSERT and UPDATE/SET operations
// It holds column name and the value to assign.
type Assignment struct {
	Column        string
	ColumnValue   interface{}
	OriginalValue interface{}
}

// Assignments is a list of Assignment
type Assignments []Assignment

// Args are the list of values of the Assignment list
func (a Assignments) Args() []interface{} {
	args := make([]interface{}, len(a))
	for i := range a {
		args[i] = a[i].ColumnValue
	}
	return args
}

// Add adds an assignment to the list
func (a *Assignments) Add(name string, columnValue interface{}, originalValue interface{}) {
	*a = append(*a, Assignment{
		Column:        name,
		ColumnValue:   columnValue,
		OriginalValue: originalValue,
	})
}
