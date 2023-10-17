package transaction

import (
	"context"
	"github.com/google/uuid"

	"github.com/pkg/errors"
	"wallet-se/internal/entity"
	"wallet-se/internal/presentations"
	"wallet-se/internal/repositories"
)

type service struct {
	repo *repositories.Repository
}

func NewService(repo *repositories.Repository) Transaction {
	return &service{repo: repo}
}

func (s *service) GetAllTransaction(ctx context.Context) ([]entity.Transaction, error) {
	transactions, err := s.repo.Transaction.GetAllTransaction(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting all transactions on ")
	}

	return transactions, nil
}

func (s *service) GetTransactionByID(ctx context.Context, transactionID string) (*entity.Transaction, error) {
	transactions, err := s.repo.Transaction.FindTransactionByID(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting transaction id %s", transactionID)
	}

	return transactions, nil
}

func (s *service) CreateTransaction(ctx context.Context, input presentations.TransactionCreate) (*entity.Transaction, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation(s) error")
	}

	//id := uuid.NewString()
	input.ID = uuid.NewString()
	err := s.repo.Transaction.CreateTransaction(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "creating transaction")

	}

	transactions, err := s.repo.Transaction.FindTransactionByID(ctx, input.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting transaction id %s", input.ID)
	}

	return transactions, nil
}
