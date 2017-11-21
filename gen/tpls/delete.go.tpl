package {{.Package}}

// String returns the SQL DELETE statement string
func (s *TDelete) String() string {
	return "DELETE FROM '{{.Type.Table}}' " + s.where.String()
}

