package orm

// Dialect is
type Dialect interface {
	Name() string
	Create() string
}

type Columner interface {
	Columns() []string
	Count() bool
}

// Where is an API for creating WHERE statements
type Where interface {
	StatementArger
	// Or applies OR condition with another Where statement
	Or(Where) Where
	// And applies AND condition with another Where statement
	And(Where) Where
}

type StatementArger interface {
	Statement() string
	Args() []interface{}
}

// Op is an SQL comparison operation
type Op string

const (
	OpEq   = "="
	OpNe   = "<>"
	OpGt   = ">"
	OpGE   = ">="
	OpLt   = "<"
	OpLE   = "<="
	OpLike = "LIKE"
)

type OrderDir string

const (
	Asc  OrderDir = "ASC"
	Desc OrderDir = "DESC"
)

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
