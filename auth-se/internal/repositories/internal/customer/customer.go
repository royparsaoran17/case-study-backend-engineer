package customer

import "auth-se/pkg/postgres"

type customer struct {
	db postgres.Adapter
}

func NewCustomer(db postgres.Adapter) *customer {
	return &customer{
		db: db,
	}
}

func (c *customer) Sortable(field string) bool {
	switch field {
	case "created_at", "updated_at", "name", "phone":
		return true
	default:
		return false
	}

}

func (c *customer) Searchable(field string) bool {
	switch field {
	case "name", "role_id", "phone":
		return true
	default:
		return false
	}

}
