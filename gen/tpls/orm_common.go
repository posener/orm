package tpls

import (
	"database/sql"

	"github.com/posener/orm"
)

type ORM struct {
	dialect orm.Dialect
	db      *sql.DB
	logger  Logger
}

func (o *ORM) Close() error {
	return o.db.Close()
}

// Create returns a struct for a CREATE statement
func (o *ORM) Create() *Create {
	return &Create{orm: o}
}

// Logger sets a logger to the ORM package
func (o *ORM) Logger(logger Logger) {
	o.logger = logger
}
