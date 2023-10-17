package repositories

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"wallet-se/internal/repositories/internal/transaction"
	"wallet-se/internal/repositories/internal/wallet"
	"wallet-se/pkg/databasex"
)

type Repository struct {
	db          databasex.Adapter
	Transaction Transaction
	Wallet      Wallet
}

func NewRepository(db databasex.Adapter) *Repository {
	return &Repository{
		db:          db,
		Transaction: transaction.NewTransaction(db),
		Wallet:      wallet.NewWallet(db),
	}
}

func (r Repository) BeginTx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error) {
	return r.db.BeginTx(ctx, options)
}

func (r Repository) CommitTx(ctx context.Context, tx *sqlx.Tx) error {
	return r.db.CommitTx(ctx, tx)
}

func (r Repository) RollbackTx(ctx context.Context, tx *sqlx.Tx) error {
	return r.db.RollbackTx(ctx, tx)
}
