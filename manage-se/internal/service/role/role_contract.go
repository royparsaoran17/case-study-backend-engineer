package role

import (
	"context"
	"manage-se/internal/provider/auth"
)

type Role interface {
	GetAllRole(ctx context.Context) ([]auth.Role, error)
}
