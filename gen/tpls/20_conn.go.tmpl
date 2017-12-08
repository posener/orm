// {{.Type.PrefixPrivate}}Conn represents a DB connection for manipulating a given struct.
// All functions available to interact with an SQL table that is related
// to this struct, are done by an instance of this object.
// To get an instance of orm use Open or New functions.
type {{.Type.PrefixPrivate}}Conn struct {
	dialect dialect.Dialect
	db      orm.DB
	logger  orm.Logger
}

func (c *{{.Type.PrefixPrivate}}Conn) Close() error {
	return c.db.Close()
}

// Logger sets a logger to the conn package
func (c *{{.Type.PrefixPrivate}}Conn) Logger(logger orm.Logger) {
	c.logger = logger
}

// {{$.Type.PrefixPublic}}CreateBuilder builds an SQL CREATE statement parameters
type {{$.Type.PrefixPublic}}CreateBuilder struct {
	params common.CreateParams
	conn   *{{.Type.PrefixPrivate}}Conn
}

// IfNotExists sets IF NOT EXISTS for the CREATE SQL statement
func (b *{{$.Type.PrefixPublic}}CreateBuilder) IfNotExists() *{{$.Type.PrefixPublic}}CreateBuilder {
	b.params.IfNotExists = true
	return b
}

// Context sets the context for the SQL query
func (b *{{$.Type.PrefixPublic}}CreateBuilder) Context(ctx context.Context) *{{$.Type.PrefixPublic}}CreateBuilder {
	b.params.Ctx = ctx
	return b
}

// {{$.Type.PrefixPublic}}InsertBuilder builds an INSERT statement parameters
type {{$.Type.PrefixPublic}}InsertBuilder struct {
	params common.InsertParams
	conn   *{{.Type.PrefixPrivate}}Conn
}

// Context sets the context for the SQL query
func (b *{{$.Type.PrefixPublic}}InsertBuilder) Context(ctx context.Context) *{{$.Type.PrefixPublic}}InsertBuilder {
	b.params.Ctx = ctx
	return b
}

// {{$.Type.PrefixPublic}}UpdateBuilder builds SQL INSERT statement parameters
type {{$.Type.PrefixPublic}}UpdateBuilder struct {
	params common.UpdateParams
	conn   *{{.Type.PrefixPrivate}}Conn
}

// Where sets the WHERE statement to the SQL query
func (b *{{$.Type.PrefixPublic}}UpdateBuilder) Where(where common.Where) *{{$.Type.PrefixPublic}}UpdateBuilder {
	b.params.Where = where
	return b
}

// Context sets the context for the SQL query
func (b *{{$.Type.PrefixPublic}}UpdateBuilder) Context(ctx context.Context) *{{$.Type.PrefixPublic}}UpdateBuilder {
	b.params.Ctx = ctx
	return b
}

// {{$.Type.PrefixPublic}}DeleteBuilder builds SQL DELETE statement parameters
type {{$.Type.PrefixPublic}}DeleteBuilder struct {
	params common.DeleteParams
	conn   *{{.Type.PrefixPrivate}}Conn
}

// Where applies where conditions on the SQL query
func (b *{{$.Type.PrefixPublic}}DeleteBuilder) Where(w common.Where) *{{$.Type.PrefixPublic}}DeleteBuilder {
	b.params.Where = w
	return b
}

// Context sets the context for the SQL query
func (b *{{$.Type.PrefixPublic}}DeleteBuilder) Context(ctx context.Context) *{{$.Type.PrefixPublic}}DeleteBuilder {
	b.params.Ctx = ctx
	return b
}

// log if a logger was set
func (c *{{.Type.PrefixPrivate}}Conn) log(s string, args ...interface{}) {
	if c.logger == nil {
		return
	}
	c.logger(s, args...)
}
