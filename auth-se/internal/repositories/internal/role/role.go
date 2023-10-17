package role

import "auth-se/pkg/postgres"

type role struct {
	db postgres.Adapter
}

func NewRole(db postgres.Adapter) *role {
	return &role{
		db: db,
	}
}
