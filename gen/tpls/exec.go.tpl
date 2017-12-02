package {{.Package}}

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"unsafe"
	"github.com/posener/orm"

    "{{.Type.ImportPath}}"
)

// Exec creates a table for the given struct
func (b *CreateBuilder) Exec() (sql.Result, error) {
	stmt, args := b.conn.dialect.Create(&b.params)
	b.conn.log("Create: '%v' %v", stmt, args)
	return b.conn.db.ExecContext(contextOrBackground(b.params.Ctx), stmt, args...)
}

// query is used by the Select.Query and Select.Limit functions
func (b *SelectBuilder) query(ctx context.Context) (*sql.Rows, error) {
	stmt, args := b.conn.dialect.Select(&b.params)
	b.conn.log("Query: '%v' %v", stmt, args)
	return b.conn.db.QueryContext(ctx, stmt, args...)
}

// Exec inserts the data to the given database
func (b *InsertBuilder) Exec() (sql.Result, error) {
	if len(b.params.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to insert")
	}
	stmt, args := b.conn.dialect.Insert(&b.params)
	b.conn.log("Insert: '%v' %v", stmt, args)
	return b.conn.db.ExecContext(contextOrBackground(b.params.Ctx), stmt, args...)
}

// Exec inserts the data to the given database
func (b *UpdateBuilder) Exec() (sql.Result, error) {
	if len(b.params.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	stmt, args := b.conn.dialect.Update(&b.params)
	b.conn.log("Update: '%v' %v", stmt, args)
	return b.conn.db.ExecContext(contextOrBackground(b.params.Ctx), stmt, args...)
}

// Exec runs the delete statement on a given database.
func (b *DeleteBuilder) Exec() (sql.Result, error) {
	stmt, args := b.conn.dialect.Delete(&b.params)
	b.conn.log("Delete: '%v' %v", stmt, args)
	return b.conn.db.ExecContext(contextOrBackground(b.params.Ctx), stmt, args...)
}

// Query the database
func (b *SelectBuilder) Query() ([]{{.Type.ExtName}}, error) {
    ctx := contextOrBackground(b.params.Ctx)
    rows, err := b.query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.ExtName}}
	for rows.Next() {
	    // check context cancellation
	    if err := ctx.Err(); err != nil  {
	        return nil, err
	    }
		item, err := b.selector.First(b.conn.dialect.Name(), values(*rows))
        if err != nil {
			return nil, err
		}
		all = append(all, *item)
	}
	return all, rows.Err()
}

// Count add a count column to the query
func (b *SelectBuilder) Count() ([]{{.Type.Name}}Count, error) {
    ctx := contextOrBackground(b.params.Ctx)
    b.selector.count = true
    rows, err := b.query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.Name}}Count
	for rows.Next() {
	    // check context cancellation
	    if err := ctx.Err(); err != nil  {
	        return nil, err
	    }
		item, err := b.selector.FirstCount(b.conn.dialect.Name(), values(*rows))
        if err != nil {
			return nil, err
		}
		all = append(all, *item)
	}
	return all, rows.Err()
}

// First returns the first row that matches the query.
// If no row matches the query, an ErrNotFound will be returned.
// This call cancels any paging that was set with the
// SelectBuilder previously.
func (b *SelectBuilder) First() (*{{.Type.ExtName}}, error) {
    ctx := contextOrBackground(b.params.Ctx)
    b.params.Page.Limit = 1
    b.params.Page.Offset = 0
    rows, err := b.query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := rows.Next()
    if !found {
        return nil, orm.ErrNotFound
    }
    item, err := b.selector.First(b.conn.dialect.Name(), values(*rows))
    if err != nil {
        return nil, err
    }
	return item, rows.Err()
}

func contextOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
	    return context.Background()
	}
	return ctx
}

// values is a hack to the sql.Rows struct
// since the rows struct does not expose it's lastcols values, or a way to give
// a custom scanner to the Scan method.
// See issue https://github.com/golang/go/issues/22544
func values(r sql.Rows) []driver.Value {
	// some ugly hack to access lastcols field
	rs := reflect.ValueOf(&r).Elem()
	rf := rs.FieldByName("lastcols")

	// overcome panic reflect.Value.Interface: cannot return value obtained from unexported field or method
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	return rf.Interface().([]driver.Value)
}
