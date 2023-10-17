package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/lib/pq"
)

// Error is a pgSQL database error.
type Error string

// Error implement error interface.
func (e Error) Error() string {
	return string(e)
}

// Err are known errors for Postgres SQL.
const (
	ErrUniqueViolation     = Error("unique_violation")
	ErrNullValueNotAllowed = Error("null_value_not_allowed")
	ErrorUndefinedTable    = Error("undefined_table")
)

// canceledMessage is an error that occurs when deadline exceeded.
const canceledMessage = "pq: canceling statement due to user request"

// ParseSQLError parses driver specific error into known errors.
// this function converts the error from pq driver using postgres SQL error codes.
// https://www.postgresql.org/docs/current/errcodes-appendix.html
func (d *adapter) ParseSQLError(err error) error {
	// Parse by value
	switch err {
	case sql.ErrNoRows:
		return sql.ErrNoRows
	case driver.ErrBadConn:
		return context.DeadlineExceeded
	}

	// Parse by type
	switch et := err.(type) {
	case *pq.Error:
		switch et.Code {
		case "23505":
			return ErrUniqueViolation
		case "42P01":
			return ErrorUndefinedTable
		case "22004":
			return ErrNullValueNotAllowed
		}
	}

	// Parse by message
	switch err.Error() {
	case canceledMessage:
		return context.DeadlineExceeded
	}

	return err
}
