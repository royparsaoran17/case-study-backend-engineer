package wallet

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
	"wallet-se/internal/common"
	"wallet-se/internal/consts"
	"wallet-se/internal/entity"
	"wallet-se/internal/presentations"
	"wallet-se/internal/repositories"
	"wallet-se/internal/repositories/repooption"
)

type service struct {
	repo *repositories.Repository
}

func NewService(repo *repositories.Repository) Wallet {
	return &service{repo: repo}
}

func (s *service) GetAllWallet(ctx context.Context, meta *common.Metadata) ([]entity.Wallet, error) {
	wallets, err := s.repo.Wallet.GetAllWallet(ctx, meta)
	if err != nil {
		return nil, errors.Wrap(err, "getting all wallets on ")
	}

	return wallets, nil
}

func (s *service) GetWalletByID(ctx context.Context, walletID string) (*entity.WalletDetail, error) {
	wallets, err := s.repo.Wallet.FindWalletByID(ctx, walletID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting wallet id %s", walletID)
	}

	return wallets, nil
}

func (s *service) GetWalletByOwned(ctx context.Context, ownedBy string) (*entity.WalletDetail, error) {
	wallets, err := s.repo.Wallet.FindWalletByOwned(ctx, ownedBy)
	if err != nil {
		return nil, errors.Wrapf(err, "getting owned id %s", ownedBy)
	}

	return wallets, nil
}

func (s *service) UpdateWalletByID(ctx context.Context, walletID string, input presentations.WalletUpdate) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation(s) error")
	}

	_, err := s.repo.Wallet.FindWalletByID(ctx, walletID)
	if err != nil {
		return errors.Wrapf(err, "getting wallet id %s", walletID)
	}

	if err := s.repo.Wallet.UpdateWallet(ctx, walletID, input); err != nil {
		return errors.Wrap(err, "updating wallet")

	}

	return nil
}

func (s *service) CreateWallet(ctx context.Context, input presentations.WalletCreate) (*entity.WalletDetail, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation(s) error")
	}

	walletID := uuid.NewString()
	input.ID = walletID
	err := s.repo.Wallet.CreateWallet(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "creating wallet")
	}

	wallet, err := s.repo.Wallet.FindWalletByID(ctx, walletID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting wallet id %s", walletID)
	}

	return wallet, nil
}

func (s *service) WalletDeposit(ctx context.Context, input presentations.WalletDeposit) (*entity.WalletDetail, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation(s) error")
	}

	wallet, err := s.repo.Wallet.FindWalletByID(ctx, input.WalletID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	if wallet.Status == consts.StatusWalletDisabled.String() {
		return nil, consts.ErrWalletDisabled
	}

	tx, err := s.repo.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		return nil, errors.Wrap(err, "begin tx")
	}

	err = s.repo.Transaction.CreateTransaction(ctx, presentations.TransactionCreate{
		ID:            uuid.NewString(),
		Type:          consts.TypeDebit.String(),
		Amount:        input.Balance,
		Status:        "success",
		TransactionAt: time.Now().Local(),
		TransactionBy: input.CustomerID,
		WalletID:      wallet.ID,
	}, repooption.WithTx(tx))
	if err != nil {
		errRollBack := s.repo.RollbackTx(ctx, tx)
		if errRollBack != nil {
			err = errors.Wrap(err, errRollBack.Error())
		}

		return nil, errors.Wrap(err, "create transaction")
	}

	err = s.repo.Wallet.UpdateWallet(ctx, wallet.ID, presentations.WalletUpdate{
		Balance: wallet.Balance + input.Balance,
		Status:  wallet.Status,
	}, repooption.WithTx(tx))
	if err != nil {
		errRollBack := s.repo.RollbackTx(ctx, tx)
		if errRollBack != nil {
			err = errors.Wrap(err, errRollBack.Error())
		}

		return nil, errors.Wrap(err, "update wallet wallet")
	}

	errCommit := s.repo.CommitTx(ctx, tx)
	if errCommit != nil {
		return nil, errors.Wrap(errCommit, "commit transaction")
	}

	updatedWallet, err := s.repo.Wallet.FindWalletByID(ctx, input.WalletID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting wallet id %s", input.WalletID)
	}

	return updatedWallet, nil
}

func (s *service) WalletWithdraw(ctx context.Context, input presentations.WalletWithdraw) (*entity.WalletDetail, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation(s) error")
	}

	wallet, err := s.repo.Wallet.FindWalletByID(ctx, input.WalletID)
	if err != nil {
		return nil, errors.Wrap(err, "get user wallet")
	}

	if wallet.Status == consts.StatusWalletDisabled.String() {
		return nil, consts.ErrWalletDisabled
	}

	tx, err := s.repo.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		return nil, errors.Wrap(err, "begin tx")
	}

	err = s.repo.Transaction.CreateTransaction(ctx, presentations.TransactionCreate{
		ID:            uuid.NewString(),
		Type:          consts.TypeCredit.String(),
		Amount:        input.Balance,
		Status:        "success",
		TransactionAt: time.Now().Local(),
		TransactionBy: input.CustomerID,
		WalletID:      wallet.ID,
	}, repooption.WithTx(tx))
	if err != nil {
		errRollBack := s.repo.RollbackTx(ctx, tx)
		if errRollBack != nil {
			err = errors.Wrap(err, errRollBack.Error())
		}

		return nil, errors.Wrap(err, "create transaction")
	}

	err = s.repo.Wallet.UpdateWallet(ctx, wallet.ID, presentations.WalletUpdate{
		Balance: wallet.Balance - input.Balance,
		Status:  wallet.Status,
	}, repooption.WithTx(tx))
	if err != nil {
		errRollBack := s.repo.RollbackTx(ctx, tx)
		if errRollBack != nil {
			err = errors.Wrap(err, errRollBack.Error())
		}

		return nil, errors.Wrap(err, "update wallet wallet")
	}

	errCommit := s.repo.CommitTx(ctx, tx)
	if errCommit != nil {
		return nil, errors.Wrap(errCommit, "commit transaction")
	}

	updatedWallet, err := s.repo.Wallet.FindWalletByID(ctx, input.WalletID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting wallet id %s", input.WalletID)
	}

	return updatedWallet, nil
}
