package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"auth-se/pkg/tracer"
	"auth-se/pkg/util"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// CreateSession create new session maria db
func CreateSession(cfg *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		connStringTemplate,
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	to := int64(cfg.Timeout.Seconds())

	if to > 0 {
		dsn = fmt.Sprintf("%s connect_timeout=%d", dsn, to)
	}

	if tz := cfg.Timezone; len(tz) > 0 {
		dsn = fmt.Sprintf("%s timezone=%s", dsn, tz)
	}

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		err = fmt.Errorf("opening db: %w", err)

		return db, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	return db, nil
}

type adapter struct {
	cfg *Config
	db  *sqlx.DB
}

// NewAdapter initialize single postgres db
func NewAdapter(cfg *Config) (Adapter, error) {
	db, err := CreateSession(cfg)
	if err != nil {
		err = fmt.Errorf("creating db session: %w", err)

		return nil, err
	}

	return &adapter{
		cfg: cfg,
		db:  db,
	}, err
}

// QueryRow select single row database will return  sql.row raw
func (d *adapter) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.query_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)

	defer tracer.SpanFinish(ctx)

	return d.db.QueryRowContext(ctx, query, args...)
}

// QueryRows select multiple rows of database will return  sql.rows raw
func (d *adapter) QueryRows(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.query_rows",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)

	defer tracer.SpanFinish(ctx)

	return d.db.QueryContext(ctx, query, args...)
}

// Fetch select multiple rows of database will cast data to struct passing by parameter
func (d *adapter) Fetch(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.fetch_rows",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)

	defer tracer.SpanFinish(ctx)

	return d.db.SelectContext(ctx, dst, query, args...)
}

// FetchRow fetching one row database will cast data to struct passing by parameter
func (d *adapter) FetchRow(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.fetch_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)

	defer tracer.SpanFinish(ctx)

	return d.db.GetContext(ctx, dst, query, args...)
}

// Exec execute mysql command query
func (d *adapter) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.exec",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)

	defer tracer.SpanFinish(ctx)

	return d.db.ExecContext(ctx, query, args...)
}

// BeginTx start new transaction session
func (d *adapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.begin_transaction")

	defer tracer.SpanFinish(ctx)

	return d.db.BeginTxx(ctx, opts)
}

// Ping check database connectivity
func (d *adapter) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// HealthCheck checking healthy of database connection
func (d *adapter) HealthCheck() error {
	return d.Ping(context.Background())
}

func (d *adapter) ExecTx(ctx context.Context, tx *sqlx.Tx, query string, args ...any) (sql.Result, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.exec_transaction",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		tracer.SpanError(ctx, err)
		return nil, err
	}

	return result, nil
}

func (d *adapter) QueryTx(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) (*sql.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.query_transaction",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		tracer.SpanError(ctx, err)
		return nil, err
	}

	return rows, nil
}

func (d *adapter) CommitTx(ctx context.Context, tx *sqlx.Tx) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.commit_transaction")
	defer tracer.SpanFinish(ctx)

	err := tx.Commit()
	if err != nil {
		tracer.SpanError(ctx, err)
		return err
	}

	return nil
}

func (d *adapter) RollbackTx(ctx context.Context, tx *sqlx.Tx) error {
	ctx = tracer.DBSpanStartWithOption(ctx, d.cfg.Name, "db.rollback_transaction")
	defer tracer.SpanFinish(ctx)

	err := tx.Rollback()
	if err != nil {
		tracer.SpanError(ctx, err)
		return err
	}

	return nil
}
