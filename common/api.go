package common

import "strings"

// Op is an SQL comparison operation
type Op string

// Selector is interface for generating columns of SELECT queries.
// With this interface, a dialect talks to struct specific generated implementation.
type Selector interface {
	Columns() []string
	Joins() []JoinParams
	Count() bool
}

// JoinParams are parameters to perform a join operation
// Field SelectParams is used to perform select operations on the join struct's field.
// Pairings describe the relation between the join's fields
type JoinParams struct {
	SelectParams
	Pairings []Pairing
}

// Pairing describe a join relation
type Pairing struct {
	// Column is the column in the current table for the JOIN statement
	Column string
	// JoinedColumn is the column in the referenced table for the JOIN statement
	JoinedColumn string
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

// StatementArger is interface for queries.
// The statement and the args are given to the SQL query.
type StatementArger interface {
	Statement(string) string
	Args() []interface{}
}

// OrderDir is direction in which a column can be ordered
type OrderDir string

// Order is a struct that is used for ORDER BY operations.
// It holds the column name and the order direction.
type Order struct {
	Column string
	Dir    OrderDir
}

// Orders is a list of Order
type Orders []Order

// Add adds an order of a column and direction
func (g *Orders) Add(name string, dir OrderDir) {
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
	Column string
	Value  interface{}
}

// Assignments is a list of Assignment
type Assignments []Assignment

// Args are the list of values of the Assignment list
func (a Assignments) Args() []interface{} {
	args := make([]interface{}, len(a))
	for i := range a {
		args[i] = a[i].Value
	}
	return args
}

// Add adds an assignment to the list
func (a *Assignments) Add(name string, value interface{}) {
	*a = append(*a, Assignment{Column: name, Value: value})
}
