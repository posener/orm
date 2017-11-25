package orm

import "github.com/posener/orm/common"

const (
	OpEq   common.Op = "="
	OpNe   common.Op = "<>"
	OpGt   common.Op = ">"
	OpGE   common.Op = ">="
	OpLt   common.Op = "<"
	OpLE   common.Op = "<="
	OpLike common.Op = "LIKE"
)

const (
	Asc  common.OrderDir = "ASC"
	Desc common.OrderDir = "DESC"
)
