package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

const connStringTemplate = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	Timeout      time.Duration
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	Timezone     string
}

type Adapter interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRows(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error)
	Fetch(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	FetchRow(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Ping(ctx context.Context) error
	HealthCheck() error

	ExecTx(ctx context.Context, tx *sqlx.Tx, query string, args ...any) (sql.Result, error)
	QueryTx(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) (*sql.Rows, error)
	CommitTx(ctx context.Context, tx *sqlx.Tx) error
	RollbackTx(ctx context.Context, tx *sqlx.Tx) error

	ParseSQLError(err error) error
}

type Transaction interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
