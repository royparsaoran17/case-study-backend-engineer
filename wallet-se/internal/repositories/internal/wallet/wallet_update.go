package wallet

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"time"
	"wallet-se/internal/consts"
	"wallet-se/internal/presentations"
	"wallet-se/internal/repositories/repooption"
	"wallet-se/pkg/databasex"
)

func (w wallet) UpdateWallet(ctx context.Context, walletID string, input presentations.WalletUpdate, opts ...repooption.TxOption) error {

	txOpt := repooption.TxOptions{
		Tx:              nil,
		NotCommitInRepo: false,
	}

	for _, opt := range opts {
		opt(&txOpt)
	}

	if txOpt.Tx == nil {
		tx, err := w.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		if err != nil {
			return errors.Wrap(err, "beginTx")
		}

		txOpt.Tx = tx
	}

	tx := txOpt.Tx
	query := `
	UPDATE wallet 
	SET 
	    balance = $2, 
	    status = $3, 
	    updated_at = $4 
	WHERE id = $1 
	AND deleted_at is null;`

	values := []interface{}{
		walletID,
		input.Balance,
		input.Status,
		time.Now().Local(),
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

	if !txOpt.NotCommitInRepo {
		err := w.db.CommitTx(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "commit add merchant")
		}
	}

	return nil
}
