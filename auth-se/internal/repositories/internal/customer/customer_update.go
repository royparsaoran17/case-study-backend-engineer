package customer

import (
	"auth-se/internal/presentations"
	"auth-se/pkg/postgres"
	"context"
	"database/sql"
	"time"

	"auth-se/internal/consts"
	"github.com/pkg/errors"
)

func (c customer) UpdateCustomer(ctx context.Context, customerID string, input presentations.CustomerUpdate) error {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return errors.Wrap(err, "failed begin tx")
	}

	query := `
	UPDATE customers 
	SET 
	    name = $2, 
	    phone = $3, 
	    role_id = $4, 
	    updated_at = $5 
	WHERE id = $1 
	AND deleted_at is null;`

	values := []interface{}{
		customerID,
		input.Name,
		input.Phone,
		input.RoleID,
		time.Now().Local(),
	}

	if _, err := c.db.ExecTx(ctx, tx, query, values...); err != nil {
		errRollback := c.db.RollbackTx(ctx, tx)
		if errRollback != nil {
			return errors.Wrap(err, "rollback failed")
		}

		errSql := c.db.ParseSQLError(err)

		if errSql != nil {
			switch errSql {
			case postgres.ErrUniqueViolation:
				return consts.ErrPhoneAlreadyExist

			default:
				return errors.Wrap(err, "failed execute query")
			}
		}

	}

	if err := c.db.CommitTx(ctx, tx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
