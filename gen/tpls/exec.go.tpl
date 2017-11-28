package {{.Package}}

import (
	"context"
	"database/sql"
	"fmt"

	api "github.com/posener/orm"

    "{{.Type.ImportPath}}"
)

// Exec creates a table for the given struct
func (b *CreateBuilder) Exec() (sql.Result, error) {
	stmt, args := b.orm.dialect.Create(&b.params)
	b.orm.log("Create: '%v' %v", stmt, args)
	return b.orm.db.ExecContext(contextOrBackground(b.params.Ctx), stmt, args...)
}

// query is used by the Select.Query and Select.Limit functions
func (b *SelectBuilder) query(ctx context.Context) (*sql.Rows, error) {
	stmt, args := b.orm.dialect.Select(&b.params)
	b.orm.log("Query: '%v' %v", stmt, args)
	return b.orm.db.QueryContext(ctx, stmt, args...)
}

// Exec inserts the data to the given database
func (b *InsertBuilder) Exec() (sql.Result, error) {
	if len(b.params.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to insert")
	}
	stmt, args := b.orm.dialect.Insert(&b.params)
	b.orm.log("Insert: '%v' %v", stmt, args)
	return b.orm.db.ExecContext(contextOrBackground(b.params.Ctx), stmt, args...)
}

// Exec inserts the data to the given database
func (b *UpdateBuilder) Exec() (sql.Result, error) {
	if len(b.params.Assignments) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}
	stmt, args := b.orm.dialect.Update(&b.params)
	b.orm.log("Update: '%v' %v", stmt, args)
	return b.orm.db.ExecContext(contextOrBackground(b.params.Ctx), stmt, args...)
}

// Exec runs the delete statement on a given database.
func (b *DeleteBuilder) Exec() (sql.Result, error) {
	stmt, args := b.orm.dialect.Delete(&b.params)
	b.orm.log("Delete: '%v' %v", stmt, args)
	return b.orm.db.ExecContext(contextOrBackground(b.params.Ctx), stmt, args...)
}

// Query the database
func (b *SelectBuilder) Query() ([]{{.Type.FullName}}, error) {
    ctx := contextOrBackground(b.params.Ctx)
    rows, err := b.query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// extract rows to structures
	var all []{{.Type.FullName}}
	for rows.Next() {
	    // check context cancellation
	    if err := ctx.Err(); err != nil  {
	        return nil, err
	    }
		item, err := scan(b.orm.dialect.Name(), b.columns, rows)
        if err != nil {
			return nil, err
		}
		all = append(all, item.{{.Type.Name}})
	}
	return all, rows.Err()
}

// Count add a count column to the query
func (b *SelectBuilder) Count() ([]{{.Type.Name}}Count, error) {
    ctx := contextOrBackground(b.params.Ctx)
    b.columns.count = true
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
		item, err := scan(b.orm.dialect.Name(), b.columns, rows)
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
func (b *SelectBuilder) First() (*{{.Type.FullName}}, error) {
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
        return nil, api.ErrNotFound
    }
    item, err := scan(b.orm.dialect.Name(), b.columns, rows)
    if err != nil {
        return nil, err
    }
	return &item.{{.Type.Name}}, rows.Err()
}

func contextOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
	    return context.Background()
	}
	return ctx
}