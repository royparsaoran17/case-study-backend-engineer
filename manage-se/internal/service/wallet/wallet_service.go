package wallet

import (
	"context"
	"github.com/pkg/errors"
	"manage-se/internal/consts"
	"manage-se/internal/entity"
	"manage-se/internal/presentations"
	"manage-se/internal/provider"
	"manage-se/internal/provider/wallet"
)

type service struct {
	provider *provider.Provider
}

func NewService(provider *provider.Provider) Wallet {
	return &service{provider: provider}
}

func (s *service) InitializeAccountWallet(ctx context.Context) (*wallet.WalletDetail, error) {

	customer, ok := ctx.Value(consts.CtxUserAuth).(entity.CustomerContext)
	if !ok {
		return nil, consts.ErrUserUnauthorized
	}

	wallet, err := s.provider.Wallet.CreateWallet(ctx, presentations.WalletCreate{
		OwnedBy: customer.ID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating user wallet")
	}

	return wallet, nil
}

func (s *service) EnableWallet(ctx context.Context) (*wallet.WalletDetail, error) {

	customer, ok := ctx.Value(consts.CtxUserAuth).(entity.CustomerContext)
	if !ok {
		return nil, consts.ErrUserUnauthorized
	}

	wallet, err := s.provider.Wallet.GetWalletByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	if wallet.Status == consts.StatusWalletEnabled.String() {
		return nil, consts.ErrWalletAlreadyEnabled
	}

	err = s.provider.Wallet.UpdateWallet(ctx, wallet.ID, presentations.WalletUpdate{
		Balance: wallet.Balance,
		Status:  consts.StatusWalletEnabled.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "enabled wallet wallet")
	}

	updatedWallet, err := s.provider.Wallet.GetWalletByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	return updatedWallet, nil
}

func (s *service) DisableWallet(ctx context.Context) (*wallet.WalletDetail, error) {

	customer, ok := ctx.Value(consts.CtxUserAuth).(entity.CustomerContext)
	if !ok {
		return nil, consts.ErrUserUnauthorized
	}

	wallet, err := s.provider.Wallet.GetWalletByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	if wallet.Status == consts.StatusWalletEnabled.String() {
		return nil, consts.ErrWalletAlreadyEnabled
	}

	err = s.provider.Wallet.UpdateWallet(ctx, wallet.ID, presentations.WalletUpdate{
		Balance: wallet.Balance,
		Status:  consts.StatusWalletDisabled.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "enabled wallet wallet")
	}

	updatedWallet, err := s.provider.Wallet.GetWalletByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	return updatedWallet, nil
}

func (s *service) ViewWalletBalance(ctx context.Context) (*entity.WalletBalance, error) {

	customer, ok := ctx.Value(consts.CtxUserAuth).(entity.CustomerContext)
	if !ok {
		return nil, consts.ErrUserUnauthorized
	}

	wallet, err := s.provider.Wallet.GetWalletByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	return &entity.WalletBalance{
		ID:        wallet.ID,
		OwnedBy:   wallet.OwnedBy,
		Customer:  customer,
		Status:    wallet.Status,
		EnabledAt: wallet.EnabledAt,
		Balance:   wallet.Balance,
	}, nil
}

func (s *service) ViewWalletTransaction(ctx context.Context) (*wallet.WalletDetail, error) {

	customer, ok := ctx.Value(consts.CtxUserAuth).(entity.CustomerContext)
	if !ok {
		return nil, consts.ErrUserUnauthorized
	}

	wallet, err := s.provider.Wallet.GetWalletByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	if wallet.Status == consts.StatusWalletDisabled.String() {
		return nil, consts.ErrWalletDisabled
	}

	return wallet, nil
}

func (s *service) Withdraw(ctx context.Context, withdraw presentations.WalletWithdraw) (*wallet.WalletDetail, error) {

	customer, ok := ctx.Value(consts.CtxUserAuth).(entity.CustomerContext)
	if !ok {
		return nil, consts.ErrUserUnauthorized
	}

	current, err := s.provider.Wallet.GetWalletByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	if current.Status == consts.StatusWalletDisabled.String() {
		return nil, consts.ErrWalletDisabled
	}

	newWalet, err := s.provider.Wallet.WithdrawWallet(ctx, wallet.WalletWithdraw{
		CustomerID: current.OwnedBy,
		WalletID:   current.ID,
		Balance:    withdraw.Balance,
	})

	return newWalet, nil
}

func (s *service) Deposit(ctx context.Context, deposit presentations.WalletDeposit) (*wallet.WalletDetail, error) {

	customer, ok := ctx.Value(consts.CtxUserAuth).(entity.CustomerContext)
	if !ok {
		return nil, consts.ErrUserUnauthorized
	}

	current, err := s.provider.Wallet.GetWalletByCustomerID(ctx, customer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	if current.Status == consts.StatusWalletDisabled.String() {
		return nil, consts.ErrWalletDisabled
	}

	newWalet, err := s.provider.Wallet.DepositWallet(ctx, wallet.WalletDeposit{
		CustomerID: current.OwnedBy,
		WalletID:   current.ID,
		Balance:    deposit.Balance,
	})

	return newWalet, nil
}
