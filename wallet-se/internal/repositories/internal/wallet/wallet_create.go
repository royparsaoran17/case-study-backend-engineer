package wallet

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"time"
	"wallet-se/internal/consts"
	"wallet-se/internal/presentations"
	"wallet-se/pkg/databasex"
)

func (w wallet) CreateWallet(ctx context.Context, input presentations.WalletCreate) error {
	tx, err := w.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return errors.Wrap(err, "failed begin tx")
	}

	query := `
	INSERT INTO 
	    wallet 
		(
        	id, 
	        balance, 
	        status, 
	        owned_by, 
	        enabled_at, 
	        created_at, 
	        updated_at
	    )
		VALUES ($1, $2, $3, $4, $5, $5, $5)`

	tNow := time.Now().Local()
	values := []interface{}{
		input.ID,
		0,
		consts.StatusWalletDisabled,
		input.OwnedBy,
		tNow,
	}

	if _, err := w.db.ExecTx(ctx, tx, query, values...); err != nil {
		errRollback := w.db.RollbackTx(ctx, tx)
		if errRollback != nil {
			return errors.Wrap(err, "rollback failed")
		}

		errSql := w.db.ParseSQLError(err)

		if errSql != nil {
			switch errSql {
			case databasex.UniqueViolation:

				return consts.ErrWalletAlreadyExist

			default:

				return errors.Wrap(err, "failed execute query")
			}
		}

	}

	if err := w.db.CommitTx(ctx, tx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
