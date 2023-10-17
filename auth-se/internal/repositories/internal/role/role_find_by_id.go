package role

import (
	"context"
	"database/sql"

	"auth-se/internal/consts"
	"auth-se/internal/entity"
	"github.com/pkg/errors"
)

func (r role) FindRoleByID(ctx context.Context, roleID string) (*entity.Role, error) {
	query := `
SELECT 
    id, 
    name, 
    created_at::timestamptz,
    updated_at::timestamptz, 
    deleted_at::timestamptz
FROM roles 
WHERE id = $1
  AND deleted_at is null
`

	var role entity.Role

	err := r.db.FetchRow(ctx, &role, query, roleID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, consts.ErrRoleNotFound
		default:
			return nil, errors.Wrap(err, "failed to fetch row from db")
		}
	}

	return &role, nil
}
