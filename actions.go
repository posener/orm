package orm

// Insert represents an SQL INSERT statement
type Insert struct {
	Table       string
	Assignments Assignments
}

// Select represents an SQL SELECT statement
type Select struct {
	Table   string
	Columns Columner
	Where   StatementArger
	Groups  Groups
	Orders  Orders
	Page    Page
}

// Delete represents an SQL DELETE statement
type Delete struct {
	Table string
	Where StatementArger
}

// Update represents an SQL UPDATE statement
type Update struct {
	Table       string
	Assignments Assignments
	Where       StatementArger
}
