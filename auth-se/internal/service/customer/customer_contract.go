package customer

import (
	"auth-se/internal/common"
	"auth-se/internal/presentations"
	"context"

	"auth-se/internal/entity"
)

type Customer interface {
	GetAllCustomer(ctx context.Context, meta *common.Metadata) ([]entity.Customer, error)
	GetCustomerByID(ctx context.Context, customerID string) (*entity.CustomerDetail, error)
	UpdateCustomerByID(ctx context.Context, customerID string, input presentations.CustomerUpdate) error
	CreateCustomer(ctx context.Context, input presentations.CustomerCreate) (*entity.CustomerDetail, error)
}
