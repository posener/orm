package tpls

// Op is an SQL operation
type Op string

const (
	OpEq      = "="
	OpNe      = "<>"
	OpGt      = ">"
	OpGE      = ">="
	OpLt      = "<"
	OpLE      = "<="
	OpBetween = "BETWEEN"
	OpLike    = "LIKE"
	OpIn      = "IN"
)
