package tpls

import (
	"github.com/posener/orm"
)

// Insert is a struct to hold information for an INSERT statement
type Insert struct {
	orm.Insert
	orm *ORM
}
