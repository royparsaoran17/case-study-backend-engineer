package transaction

import (
	"wallet-se/pkg/databasex"
)

type transaction struct {
	db databasex.Adapter
}

func NewTransaction(db databasex.Adapter) *transaction {
	return &transaction{
		db: db,
	}
}
