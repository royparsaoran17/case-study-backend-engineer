package wallet

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
	"wallet-se/internal/consts"
	"wallet-se/internal/entity"
)

func (w wallet) FindWalletByID(ctx context.Context, walletID string) (*entity.WalletDetail, error) {
	query := `SELECT 
        jsonb_build_object(
            'id', w.id,
            'balance', w.balance,
            'status', w.status,
            'owned_by', w.owned_by,
            'enabled_at', w.enabled_at::timestamptz,
            'transactions',(
                SELECT
					json_agg (
					    json_build_object(
						'id', t.id,
						'type', t.type,
						'amount', t.amount,
						'reference_id', t.reference_id,
						'status', t.status,
						'transaction_at', t.transaction_at::timestamptz,
						'transaction_by', t.transaction_by,
						'wallet_id', t.wallet_id,
						'created_at', t.created_at::timestamptz,
						'updated_at', t.updated_at::timestamptz,
						'deleted_at', t.deleted_at::timestamptz
					)
					)
                FROM transactions t
                    WHERE w.id = t.wallet_id
                    AND t.deleted_at is null
                ),
            'created_at', w.created_at::timestamptz,
            'updated_at', w.updated_at::timestamptz,
            'deleted_at', w.deleted_at::timestamptz
        )
    FROM
        wallet w
    WHERE w.id = $1
        AND w.deleted_at is null;`

	var b []byte
	err := w.db.QueryRow(ctx, &b, query, walletID)
	if err != nil {
		sqlErr := w.db.ParseSQLError(err)
		switch sqlErr {
		case sql.ErrNoRows:
			return nil, consts.ErrWalletNotFound
		default:
			return nil, errors.Wrap(err, "failed to fetch wallet from db")
		}
	}

	var role entity.WalletDetail
	if err := json.Unmarshal(b, &role); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal byte to wallet")
	}

	return &role, nil
}
