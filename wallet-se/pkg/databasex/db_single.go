package databasex

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"wallet-se/pkg/tracer"
	"wallet-se/pkg/util"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	// check in runtime implement Databaser
	_ Adapter = (*DB)(nil)
)

type DB struct {
	db *sqlx.DB
	//instanceID string
	tx   *sqlx.Tx
	conn *sqlx.Conn // the Conn of the Tx, when tx != nil
	//opts       sql.TxOptions // valid when tx != nil
	reaMode bool
	dbName  string
}

func New(db *sqlx.DB, readMode bool, sbName string) *DB {
	return &DB{
		db:      db,
		reaMode: readMode,
		dbName:  sbName,
	}
}

func (db *DB) Ping() error {
	return db.db.Ping()
}

func (db *DB) InTransaction() bool {
	return db.tx != nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.db.Close()
}

// Exec executes a SQL statement and returns the number of rows it affected.
func (db *DB) Exec(ctx context.Context, query string, args ...any) (_ int64, err error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "exec",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	if db.reaMode {
		return 0, fmt.Errorf("database mode read only")
	}

	res, err := db.execResult(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %v", err)
	}

	return n, nil
}

// execResult executes a SQL statement and returns a sql.Result.
func (db *DB) execResult(ctx context.Context, query string, args ...any) (res sql.Result, err error) {
	if db.tx != nil {
		return db.tx.ExecContext(ctx, query, args...)
	}

	return db.db.ExecContext(ctx, query, args...)
}

// Query runs the DB query.
func (db *DB) Query(ctx context.Context, dst any, query string, args ...any) error {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "query",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	if db.tx != nil {
		return db.tx.SelectContext(ctx, dst, query, args...)
	}

	return db.db.SelectContext(ctx, dst, query, args...)
}

// QueryRow runs the query and returns a single row.
func (db *DB) QueryRow(ctx context.Context, dst interface{}, query string, args ...any) error {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "query_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)

	if db.tx != nil {
		return db.tx.GetContext(ctx, dst, query, args...)
	}

	return db.db.GetContext(ctx, dst, query, args...)
}

// QueryX runs the DB query.
func (db *DB) QueryX(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "queryx",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	if db.tx != nil {
		return db.tx.QueryContext(ctx, query, args...)
	}

	return db.db.QueryContext(ctx, query, args...)
}

// QueryRowX runs the query and returns a single row.
func (db *DB) QueryRowX(ctx context.Context, query string, args ...any) *sql.Row {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "query_rowx",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
		tracer.WithOptions("sql.args", util.DumpToString(args)),
	)
	defer tracer.SpanFinish(ctx)
	if db.tx != nil {
		return db.tx.QueryRowContext(ctx, query, args...)
	}

	return db.db.QueryRowContext(ctx, query, args...)
}

// Transact executes the given function in the context of a SQL transaction at
// the given isolation level
func (db *DB) Transact(ctx context.Context, iso sql.IsolationLevel, txFunc func(*DB) error) (err error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "transaction")
	defer tracer.SpanFinish(ctx)
	if db.reaMode {
		return fmt.Errorf("database mode read only")
	}

	// For the levels which require retry, see
	// https://www.postgresql.org/docs/11/transaction-iso.html.
	opts := &sql.TxOptions{Isolation: iso}

	return db.transact(ctx, opts, txFunc)
}

func (db *DB) transact(ctx context.Context, opts *sql.TxOptions, txFunc func(*DB) error) (err error) {
	if db.InTransaction() {
		return errors.New("db transact function was called on a DB already in a transaction")
	}

	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	tx, err := conn.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	//defer func() {
	//	if p := recover(); p != nil {
	//		tx.Rollback()
	//	} else if err != nil {
	//		tx.Rollback()
	//	} else {
	//		if txErr := tx.Commit(); txErr != nil {
	//			err = fmt.Errorf("tx commit: %w", txErr)
	//		}
	//	}
	//}()

	dbtx := New(db.db, false, db.dbName)
	dbtx.tx = tx
	dbtx.conn = conn
	//dbtx.opts = *opts

	if err := txFunc(dbtx); err != nil {
		tx.Rollback()
		return fmt.Errorf("fn(tx): %w", err)
	}

	return tx.Commit()
}

// BeginTx start new transaction session
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "begin.transaction")

	defer tracer.SpanFinish(ctx)

	return db.db.BeginTxx(ctx, opts)
}

func (db *DB) ExecTx(ctx context.Context, tx *sqlx.Tx, query string, args ...any) (sql.Result, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "db.exec_transaction",
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

func (db *DB) QueryTx(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) (*sql.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "db.query_transaction",
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

func (db *DB) CommitTx(ctx context.Context, tx *sqlx.Tx) error {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "db.commit_transaction")
	defer tracer.SpanFinish(ctx)

	err := tx.Commit()
	if err != nil {
		tracer.SpanError(ctx, err)
		return err
	}

	return nil
}

func (db *DB) RollbackTx(ctx context.Context, tx *sqlx.Tx) error {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "db.rollback_transaction")
	defer tracer.SpanFinish(ctx)

	err := tx.Rollback()
	if err != nil {
		tracer.SpanError(ctx, err)
		return err
	}

	return nil
}

func (db *DB) ParseSQLError(err error) error {
	const canceledMessage = "pq: canceling statement due to user request"

	switch err {
	case sql.ErrNoRows:
		return NoRowsFound
	case driver.ErrBadConn:
		return context.DeadlineExceeded
	}

	switch et := err.(type) {
	case *pq.Error:
		switch et.Code {
		case "02000":
			return NoRowsFound
		case "23503":
			return ForeignKeyViolation
		case "23505":
			return UniqueViolation
		case "42P01":
			return UndefinedTable
		case "22004":
			return NullValueNotAllowed
		}
	}

	switch err.Error() {
	case canceledMessage:
		return context.DeadlineExceeded
	}

	return err
}
