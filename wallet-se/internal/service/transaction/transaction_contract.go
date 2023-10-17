package transaction

import (
	"context"
	"wallet-se/internal/presentations"

	"wallet-se/internal/entity"
)

type Transaction interface {
	GetAllTransaction(ctx context.Context) ([]entity.Transaction, error)
	GetTransactionByID(ctx context.Context, transactionID string) (*entity.Transaction, error)
	CreateTransaction(ctx context.Context, input presentations.TransactionCreate) (*entity.Transaction, error)
}
