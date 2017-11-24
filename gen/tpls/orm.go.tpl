package {{.Package}}

import (
    "database/sql"

	"github.com/posener/orm"
)

const table = "{{.Type.Table}}"

// Open opens database connection
func Open(dataSourceName string) (*ORM, error) {
	db, err := sql.Open("{{.Dialect.Name}}", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &ORM{db: db}, nil
}

// Select returns an object to create a SELECT statement
func (o *ORM) Select() *Select {
	s := &Select{
		Select: orm.Select{
			Table: table,
		},
		orm: o,
	}
    s.Select.Columns = &s.columns
    return s
}

// Insert returns a new INSERT statement
func (o *ORM) Insert() *Insert {
	return &Insert{
		Insert: orm.Insert{Table: table},
		orm: o,
	}
}

// Update returns a new UPDATE statement
func (o *ORM) Update() *Update {
	return &Update{
		Update: orm.Update{Table: table},
		orm: o,
    }
}

// Delete returns an object for a DELETE statement
func (o *ORM) Delete() *Delete {
	return &Delete{
		Delete: orm.Delete{Table: table},
		orm: o,
    }
}

