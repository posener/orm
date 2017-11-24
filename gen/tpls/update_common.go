package tpls

import (
	"github.com/posener/orm"
)

// Update is a struct to hold information for an INSERT statement
type Update struct {
	orm.Update
	orm *ORM
}

func (u *Update) Where(where orm.Where) *Update {
	u.Update.Where = where
	return u
}
