package customer

import (
	"auth-se/internal/common"
	"context"
	"database/sql"
	"time"

	"auth-se/internal/consts"
	"auth-se/pkg/postgres"

	"auth-se/internal/presentations"
	"github.com/pkg/errors"
)

func (c customer) CreateCustomer(ctx context.Context, input presentations.CustomerCreate) error {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return errors.Wrap(err, "failed begin tx")
	}

	tNow := time.Now().Local()

	query := `INSERT INTO customers (id, name, phone, password, role_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $6)`

	values := []interface{}{
		input.ID,
		input.Name,
		input.Phone,
		common.Hasher(input.Password),
		input.RoleID,
		tNow,
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
