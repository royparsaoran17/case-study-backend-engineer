package databasex

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Adapter interface {
	Ping() error
	InTransaction() bool
	Close() error
	Query(ctx context.Context, dst any, query string, args ...any) error
	QueryRow(ctx context.Context, dst any, query string, args ...any) error
	QueryX(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowX(ctx context.Context, query string, args ...any) *sql.Row
	Exec(ctx context.Context, query string, args ...any) (_ int64, err error)
	Transact(ctx context.Context, iso sql.IsolationLevel, txFunc func(*DB) error) (err error)

	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)

	ExecTx(ctx context.Context, tx *sqlx.Tx, query string, args ...any) (sql.Result, error)
	QueryTx(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) (*sql.Rows, error)
	CommitTx(ctx context.Context, tx *sqlx.Tx) error
	RollbackTx(ctx context.Context, tx *sqlx.Tx) error

	ParseSQLError(err error) error
}
