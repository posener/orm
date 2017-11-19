package tpls

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
