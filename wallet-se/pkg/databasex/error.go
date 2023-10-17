package databasex

// Error is a pgSQL database error.
type Error string

// Error implement error interface.
func (e Error) Error() string {
	return string(e)
}

const (
	NoRowsFound         = Error("no rows found")
	ForeignKeyViolation = Error("foreign key violation")
	UniqueViolation     = Error("unique violation")
	UndefinedTable      = Error("undefined table")
	NullValueNotAllowed = Error("null value not allowed")
)
