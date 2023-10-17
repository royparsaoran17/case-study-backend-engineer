package customer

import (
	"auth-se/internal/consts"
	"auth-se/internal/entity"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
)

func (c customer) FindCustomerByPhone(ctx context.Context, phone string) (*entity.CustomerDetail, error) {
	query := `SELECT 
        jsonb_build_object(
            'id', c.id,
            'name', c.name,
            'phone', c.phone,
            'password', c.password,
            'role_id', c.role_id,
            'role',(
                SELECT
					json_build_object(
						'id', r.id,
						'name', r.name,
						'created_at', r.created_at::timestamptz,
						'updated_at', r.updated_at::timestamptz,
						'deleted_at', r.deleted_at::timestamptz
					)
                FROM roles r
                    WHERE c.role_id = r.id
                    AND r.deleted_at is null
                ),
            'created_at', c.created_at::timestamptz,
            'updated_at', c.updated_at::timestamptz,
            'deleted_at', c.deleted_at::timestamptz
        )
    FROM
        customers c
    WHERE c.phone = $1
        AND c.deleted_at is null;`

	var b []byte
	err := c.db.QueryRow(ctx, query, phone).Scan(&b)
	if err != nil {
		sqlErr := c.db.ParseSQLError(err)
		switch sqlErr {
		case sql.ErrNoRows:
			return nil, consts.ErrCustomerNotFound
		default:
			return nil, errors.Wrap(err, "failed to fetch customer from db")
		}
	}

	var role entity.CustomerDetail
	if err := json.Unmarshal(b, &role); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal byte to customer")
	}

	return &role, nil
}
