package provider

import (
	"context"
	"manage-se/internal/common"
	"manage-se/internal/presentations"
	"manage-se/internal/provider/auth"
	"manage-se/internal/provider/wallet"
)

type Auth interface {
	Login(ctx context.Context, input presentations.Login) (*auth.CustomerDetailToken, error)
	Verify(ctx context.Context, input presentations.Verify) (*auth.CustomerDetail, error)

	CreateCustomer(ctx context.Context, input presentations.CustomerCreate) (*auth.CustomerDetail, error)
	GetListCustomers(ctx context.Context, meta *common.Metadata) ([]auth.Customer, error)

	GetListRoles(ctx context.Context) ([]auth.Role, error)
	CreateRole(ctx context.Context, input presentations.RoleCreate) (*auth.Role, error)
}

type Wallet interface {
	CreateWallet(ctx context.Context, input presentations.WalletCreate) (*wallet.WalletDetail, error)
	UpdateWallet(ctx context.Context, walletID string, input presentations.WalletUpdate) error
	GetListWallets(ctx context.Context, meta *common.Metadata) ([]wallet.Wallet, error)
	GetWalletByCustomerID(ctx context.Context, customerID string) (*wallet.WalletDetail, error)
	DepositWallet(ctx context.Context, input wallet.WalletDeposit) (*wallet.WalletDetail, error)
	WithdrawWallet(ctx context.Context, input wallet.WalletWithdraw) (*wallet.WalletDetail, error)

	GetListTransactions(ctx context.Context) ([]wallet.Transaction, error)
	CreateTransaction(ctx context.Context, input presentations.TransactionCreate) (*wallet.Transaction, error)
}
