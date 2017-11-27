package orm

import "github.com/posener/orm/common"

// Operators for SQL WHERE statements
const (
	OpEq   common.Op = "="
	OpNe   common.Op = "<>"
	OpGt   common.Op = ">"
	OpGE   common.Op = ">="
	OpLt   common.Op = "<"
	OpLE   common.Op = "<="
	OpLike common.Op = "LIKE"
)

// Directions for SQL ORDER BY statements
const (
	Asc  common.OrderDir = "ASC"
	Desc common.OrderDir = "DESC"
)
