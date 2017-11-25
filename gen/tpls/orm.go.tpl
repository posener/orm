package {{.Package}}

import (
    "database/sql"

	"github.com/posener/orm/common"
)

const table = "{{.Type.Table}}"

// Open opens database connection
func Open(dataSourceName string) (API, error) {
	db, err := sql.Open("{{.Dialect.Name}}", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &orm{db: db}, nil
}

// New returns an orm object from a db instance
func New(db DB) API {
    return &orm{db: db}
}

// Select returns an object to create a SELECT statement
func (o *orm) Select() *Select {
	s := &Select{
		internal: common.Select{Table: table},
		orm: o,
	}
    s.internal.Columns = &s.columns
    return s
}

// Insert returns a new INSERT statement
func (o *orm) Insert() *Insert {
	return &Insert{
		internal: common.Insert{Table: table},
		orm: o,
	}
}

// Update returns a new UPDATE statement
func (o *orm) Update() *Update {
	return &Update{
		internal: common.Update{Table: table},
		orm: o,
    }
}

// Delete returns an object for a DELETE statement
func (o *orm) Delete() *Delete {
	return &Delete{
		internal: common.Delete{Table: table},
		orm: o,
    }
}

