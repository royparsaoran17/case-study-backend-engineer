package role

import (
	"auth-se/internal/presentations"
	"auth-se/pkg/postgres"
	"context"
	"database/sql"
	"time"

	"auth-se/internal/consts"
	"github.com/pkg/errors"
)

func (r role) UpdateRole(ctx context.Context, roleID string, input presentations.RoleUpdate) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return errors.Wrap(err, "failed begin tx")
	}

	query := `
	UPDATE roles 
	SET 
	    name = $2, 
	    updated_at = $3 
	WHERE id = $1 
	AND deleted_at is null;`

	values := []interface{}{
		roleID,
		input.Name,
		time.Now().Local(),
	}

	if _, err := r.db.ExecTx(ctx, tx, query, values...); err != nil {
		errRollback := r.db.RollbackTx(ctx, tx)
		if errRollback != nil {
			return errors.Wrap(err, "rollback failed")
		}

		errSql := r.db.ParseSQLError(err)

		if errSql != nil {
			switch errSql {
			case postgres.ErrUniqueViolation:
				return consts.ErrRoleAlreadyExist

			default:
				return errors.Wrap(err, "failed execute query")
			}
		}

	}

	if err := r.db.CommitTx(ctx, tx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
