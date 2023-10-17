package transaction

import (
	"context"
	"database/sql"
	"time"
	"wallet-se/internal/repositories/repooption"

	"wallet-se/internal/consts"
	"wallet-se/pkg/postgres"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"wallet-se/internal/presentations"
)

func (t transaction) CreateTransaction(ctx context.Context, input presentations.TransactionCreate, opts ...repooption.TxOption) error {
	txOpt := repooption.TxOptions{
		Tx:              nil,
		NotCommitInRepo: false,
	}

	for _, opt := range opts {
		opt(&txOpt)
	}

	if txOpt.Tx == nil {
		tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		if err != nil {
			return errors.Wrap(err, "beginTx")
		}

		txOpt.Tx = tx
	}

	tx := txOpt.Tx

	tNow := time.Now().Local()

	query := `
	INSERT INTO 
	    transactions (
        	id, 
	        type,
	        amount, 
	        reference_id, 
	        status, 
	        transaction_at, 
	        transaction_by, 
	        wallet_id, 
	        created_at, 
	        updated_at
	    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)`

	values := []interface{}{
		input.ID,
		input.Type,
		input.Amount,
		uuid.New(),
		input.Status,
		input.TransactionAt,
		input.TransactionBy,
		input.WalletID,
		tNow,
	}

	if _, err := t.db.ExecTx(ctx, tx, query, values...); err != nil {
		errRollback := t.db.RollbackTx(ctx, tx)
		if errRollback != nil {
			return errors.Wrap(err, "rollback failed")
		}

		errSql := t.db.ParseSQLError(err)

		if errSql != nil {
			switch errSql {
			case postgres.ErrUniqueViolation:
				return consts.ErrTransactionAlreadyExist

			default:
				return errors.Wrap(err, "failed execute query")
			}
		}

	}

	if !txOpt.NotCommitInRepo {
		err := t.db.CommitTx(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "commit add merchant")
		}
	}

	return nil
}
