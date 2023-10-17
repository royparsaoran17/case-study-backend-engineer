package wallet

import (
	"context"
	"manage-se/internal/entity"
	"manage-se/internal/presentations"
	"manage-se/internal/provider/wallet"
)

type Wallet interface {
	InitializeAccountWallet(ctx context.Context) (*wallet.WalletDetail, error)

	Deposit(ctx context.Context, deposit presentations.WalletDeposit) (*wallet.WalletDetail, error)
	Withdraw(ctx context.Context, withdraw presentations.WalletWithdraw) (*wallet.WalletDetail, error)

	ViewWalletTransaction(ctx context.Context) (*wallet.WalletDetail, error)
	ViewWalletBalance(ctx context.Context) (*entity.WalletBalance, error)

	EnableWallet(ctx context.Context) (*wallet.WalletDetail, error)
	DisableWallet(ctx context.Context) (*wallet.WalletDetail, error)
}
