package repositories

import (
	"manage-se/pkg/postgres"
)

type Repository struct {
}

func NewRepository(db postgres.Adapter) *Repository {
	return &Repository{}
}
