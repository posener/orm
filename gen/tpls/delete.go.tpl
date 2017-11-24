package {{.Package}}

import "github.com/posener/orm/dialect/{{.Dialect.Name}}"

// String returns the SQL DELETE statement
func (s *TDelete) String() string {
    return {{.Dialect.Name}}.Delete(s.orm, s.where)
}
