// Package repositories
// Automatic generated
package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/sync/errgroup"

	"{{.ModuleName}}/internal/entity"
	"{{.ModuleName}}/internal/common"
	"{{.ModuleName}}/pkg/tracer"
	"{{.ModuleName}}/pkg/databasex"
	"{{.ModuleName}}/pkg/builderx"
)

// {{.RepoContractName}} contract of {{.ObjectName}}
type {{.RepoContractName}} interface {
    Storer
	Updater
	Deleter
	Counter
	FindOne(ctx context.Context, param any) (*entity.{{.ObjectName}}, error)
	Find(ctx context.Context, param any) ([]entity.{{.ObjectName}}, error)
	FindWithCount(ctx context.Context, param any) ([]entity.{{.ObjectName}}, int64, error)
}

type {{.StructName}} struct {
	db databasex.Adapter
}

// New{{.ObjectName}} create new instance of {{.ObjectName}}
func New{{.ObjectName}}(db databasex.Adapter) {{.RepoContractName}} {
	return &{{.StructName}}{db: db}
}

// FindOne {{.StructName}}
func (r *{{.StructName}}) FindOne(ctx context.Context, param any) (*entity.{{.ObjectName}}, error) {
	var (
		result entity.{{.ObjectName}}
		err    error
	)

	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_find_one")
	defer tracer.SpanFinish(ctx)

	wq, err := builderx.StructToMySqlQueryWhere(param, "db")
	if err != nil {
	    tracer.SpanError(ctx, err)
		return nil, err
	}

	q := `{{.RepoQuery}} %s LIMIT 1`

	err = r.db.QueryRow(ctx, &result, fmt.Sprintf(q, wq.Query), wq.Values...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

// Find {{.StructName}}
func (r *{{.StructName}}) Find(ctx context.Context, param any) ([]entity.{{.ObjectName}}, error) {
	var (
		result []entity.{{.ObjectName}}
		err    error
	)

	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_finds")
	defer tracer.SpanFinish(ctx)

	wq, err := builderx.StructToMySqlQueryWhere(param, "db")
	if err != nil {
	    tracer.SpanError(ctx, err)
		return nil, err
	}

	q := `{{.RepoQuery}} %s LIMIT ? OFFSET ? `

	vals := wq.Values
	vals = append(vals, wq.Limit, common.PageToOffset(wq.Limit, wq.Page))
	err = r.db.Query(ctx, &result, fmt.Sprintf(q, wq.Query), vals...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

// Store {{.StructName}}
func (r *{{.StructName}}) Store(ctx context.Context, param any) (int64, error) {
	var (
		err error
		affected int64
	)

	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_store")
    defer tracer.SpanFinish(ctx)

	np := &param
	param = *np
	query, vals, err := builderx.StructToQueryInsert(param, "{{.TableName}}", "db")
	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	err = r.db.Transact(ctx, sql.LevelRepeatableRead, func(tx *databasex.DB) error {
    		af, err := tx.Exec(ctx, query, vals...)
    		affected = af
    		return err
    })

    return affected, err

}

// Update {{.StructName}} data
func (r *{{.StructName}}) Update(ctx context.Context, input any, where any) (int64, error) {
	var (
		err error
		affected int64
	)

    ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_update")
    defer tracer.SpanFinish(ctx)

	query, vals, err := builderx.StructToQueryUpdate(input, where, "{{.TableName}}", "db")
	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
    err = r.db.Transact(ctx, sql.LevelRepeatableRead, func(tx *databasex.DB) error {
            af, err := tx.Exec(ctx, query, vals...)
            affected = af
            return err
    })

    return affected, err
}

// Delete {{.StructName}} from database
func (r *{{.StructName}}) Delete(ctx context.Context, param any) (int64, error) {
    var (
            err error
            affected int64
    )
    ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_delete")
	defer tracer.SpanFinish(ctx)

	query, vals, err := builderx.StructToQueryDelete(param, "{{.TableName}}", "db", true)
	if err != nil {
	    tracer.SpanError(ctx, err)
		return 0, err
	}

	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
    err = r.db.Transact(ctx, sql.LevelRepeatableRead, func(tx *databasex.DB) error {
            af, err := tx.Exec(ctx, query, vals...)
            affected = af
            return err
    })

    return affected, err
}

// Count {{.StructName}}
func (r *{{.StructName}}) Count(ctx context.Context, p any) (total int64, err error) {
	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_count")
	defer tracer.SpanFinish(ctx)

	wq, err := builderx.StructToMySqlQueryWhere(p, "db")
	if err != nil {
		tracer.SpanError(ctx, err)
		return
	}

	q := fmt.Sprintf(`
		SELECT
        	COUNT(id) AS jumlah
		FROM {{.TableName}} %s `, wq.Query)

	err = r.db.QueryRow(ctx, &total, q, wq.Values...)
	if err != nil {
		tracer.SpanError(ctx, err)
		err = err
		return
	}

	return
}

// FindWithCount find {{.StructName}} with count
func (r *{{.StructName}}) FindWithCount(ctx context.Context, param any) ([]entity.{{.ObjectName}}, int64, error) {

	var (
		cl    []entity.{{.ObjectName}}
		count int64
	)

	ctx = tracer.SpanStart(ctx, "repo.{{.FileName}}_with_count")
    defer tracer.SpanFinish(ctx)

	group, newCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		l, err := r.Find(newCtx, param)
		cl = l
		return err
	})
	group.Go(func() error {
		c, err := r.Count(ctx, param)
		count = c
		return err
	})

	err := group.Wait()

	return cl, count, err
}
