package role

import (
	"auth-se/internal/entity"
	"context"
	"github.com/pkg/errors"
)

func (r role) GetAllRole(ctx context.Context) ([]entity.Role, error) {
	query := `
SELECT 
    id, 
    name, 
    created_at::timestamptz,
    updated_at::timestamptz, 
    deleted_at::timestamptz
FROM roles 
  WHERE deleted_at is null
`

	roles := make([]entity.Role, 0)

	err := r.db.Fetch(ctx, &roles, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all roles from database")
	}

	return roles, nil
}
