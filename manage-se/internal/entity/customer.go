package entity

import "manage-se/internal/provider/auth"

type CustomerContext struct {
	ID     string    `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Phone  string    `json:"phone" db:"phone"`
	RoleID string    `json:"role_id" db:"role_id"`
	Role   auth.Role `json:"role" db:"role"`
}
