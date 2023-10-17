package transaction

import (
	"context"
	"github.com/pkg/errors"
	"wallet-se/internal/entity"
)

func (t transaction) GetAllTransaction(ctx context.Context) ([]entity.Transaction, error) {
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
  WHERE deleted_at is null
`

	data := make([]entity.Transaction, 0)

	err := t.db.Query(ctx, &data, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all transaction from database")
	}

	return data, nil
}
