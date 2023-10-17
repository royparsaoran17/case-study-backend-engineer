package transaction

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"wallet-se/internal/consts"
	"wallet-se/internal/entity"
)

func (t transaction) FindTransactionByID(ctx context.Context, transactionID string) (*entity.Transaction, error) {
	query := `
SELECT 
	id, 
	type,
	amount, 
	reference_id, 
	status, 
	transaction_at, 
	transaction_by, 
	wallet_id, 
	created_at::timestamptz, 
	updated_at::timestamptz,
	deleted_at::timestamptz
FROM transactions 
WHERE id = $1
  AND deleted_at is null
`

	var data entity.Transaction
	err := t.db.QueryRow(ctx, &data, query, transactionID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, consts.ErrTransactionNotFound
		default:
			return nil, errors.Wrap(err, "failed to fetch row from db")
		}
	}

	return &data, nil
}
