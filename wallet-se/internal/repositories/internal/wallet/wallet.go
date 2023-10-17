package wallet

import (
	"wallet-se/pkg/databasex"
)

type wallet struct {
	db databasex.Adapter
}

func NewWallet(db databasex.Adapter) *wallet {
	return &wallet{
		db: db,
	}
}

func (c *wallet) Sortable(field string) bool {
	switch field {
	case "created_at", "updated_at", "name", "phone":
		return true
	default:
		return false
	}

}

func (c *wallet) Searchable(field string) bool {
	switch field {
	case "name", "role_id", "phone":
		return true
	default:
		return false
	}

}
