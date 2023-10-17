package repositories

import (
	"context"
	"wallet-se/internal/common"
	"wallet-se/internal/entity"
	"wallet-se/internal/repositories/repooption"

	"wallet-se/internal/presentations"
)

type Transaction interface {
	CreateTransaction(ctx context.Context, input presentations.TransactionCreate, opts ...repooption.TxOption) error
	FindTransactionByID(ctx context.Context, transactionID string) (*entity.Transaction, error)
	GetAllTransaction(ctx context.Context) ([]entity.Transaction, error)
}

type Wallet interface {
	CreateWallet(ctx context.Context, input presentations.WalletCreate) error
	UpdateWallet(ctx context.Context, transactionID string, input presentations.WalletUpdate, opts ...repooption.TxOption) error
	GetAllWallet(ctx context.Context, meta *common.Metadata) ([]entity.Wallet, error)
	FindWalletByID(ctx context.Context, walletID string) (*entity.WalletDetail, error)
	FindWalletByOwned(ctx context.Context, ownedBy string) (*entity.WalletDetail, error)
}
