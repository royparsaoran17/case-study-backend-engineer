package wallet

import (
	"context"
	"wallet-se/internal/common"
	"wallet-se/internal/presentations"

	"wallet-se/internal/entity"
)

type Wallet interface {
	GetAllWallet(ctx context.Context, meta *common.Metadata) ([]entity.Wallet, error)
	GetWalletByID(ctx context.Context, walletID string) (*entity.WalletDetail, error)
	GetWalletByOwned(ctx context.Context, ownedBy string) (*entity.WalletDetail, error)
	UpdateWalletByID(ctx context.Context, walletID string, input presentations.WalletUpdate) error
	WalletDeposit(ctx context.Context, input presentations.WalletDeposit) (*entity.WalletDetail, error)
	WalletWithdraw(ctx context.Context, input presentations.WalletWithdraw) (*entity.WalletDetail, error)
	CreateWallet(ctx context.Context, input presentations.WalletCreate) (*entity.WalletDetail, error)
}
