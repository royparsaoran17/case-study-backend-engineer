package role

import (
	"auth-se/internal/presentations"
	"context"

	"auth-se/internal/entity"
)

type Role interface {
	GetAllRole(ctx context.Context) ([]entity.Role, error)
	GetRoleByID(ctx context.Context, roleID string) (*entity.Role, error)
	UpdateRoleByID(ctx context.Context, roleID string, input presentations.RoleUpdate) error
	CreateRole(ctx context.Context, input presentations.RoleCreate) error
}
