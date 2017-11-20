package tpls

import "database/sql"

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func New(db DB) *ORM {
	return &ORM{db: db}
}

type ORM struct {
	db     DB
	logger Logger
}

// Create returns a struct for a CREATE statement
func (o *ORM) Create() *TCreate {
	return &TCreate{orm: o}
}

// Select returns an object to create a SELECT statement
func (o *ORM) Select() *TSelect {
	return &TSelect{orm: o}
}

// Insert returns a new INSERT statement
func (o *ORM) Insert() *TInsert {
	return &TInsert{orm: o}
}

// Insert returns a new INSERT statement
func (o *ORM) Update() *TUpdate {
	return &TUpdate{TInsert: TInsert{orm: o}}
}

// Delete returns an object for a DELETE statement
func (o *ORM) Delete() *TDelete {
	return &TDelete{orm: o}
}

// Logger sets a logger to the ORM package
func (o *ORM) Logger(logger Logger) {
	o.logger = logger
}
