package entity

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Phone     string     `json:"phone" db:"phone"`
	Password  string     `json:"password,omitempty" db:"password"`
	RoleID    uuid.UUID  `json:"role_id" db:"role_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type CustomerDetail struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Phone     string     `json:"phone" db:"phone"`
	Password  string     `json:"password" db:"password"`
	RoleID    uuid.UUID  `json:"role_id" db:"role_id"`
	Role      Role       `json:"role" db:"role"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type CustomerDetailToken struct {
	CustomerDetail
	Token string `json:"token"`
}
