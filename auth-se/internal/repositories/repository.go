package repositories

import (
	"auth-se/internal/repositories/internal/customer"
	"auth-se/internal/repositories/internal/role"
	"auth-se/pkg/postgres"
)

type Repository struct {
	Role     Role
	Customer Customer
}

func NewRepository(db postgres.Adapter) *Repository {
	return &Repository{
		Role:     role.NewRole(db),
		Customer: customer.NewCustomer(db),
	}
}
